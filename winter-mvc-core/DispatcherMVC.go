package winter_mvc_core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type dispatcherHandler struct {
	urlByProjectPrefixMapping map[string]map[string]interface{}
	handlerInterceptorMapping map[string]HandlerInterceptorInterface
	filter HttpFilterInterface
	failure FailureResponseInterface
}

var dispatcherHandlerInstance dispatcherHandler

func init(){
	dispatcherHandlerInstance.urlByProjectPrefixMapping = make(map[string]map[string]interface{})
	dispatcherHandlerInstance.handlerInterceptorMapping = make(map[string]HandlerInterceptorInterface)
}

//	dispatcherHandler实例指针
//	获得dispatcherHandler 初始化MVC
//	初始化方式1	http.HandleFunc("/projectPrefix/", winter-mvc-core.HandlerFun())
//	初始化方式2		server := http.Server{ Handler: winter_mvc_core.GetInstanceByDispatcherHandler(), ...
func GetInstanceByDispatcherHandler() *dispatcherHandler{
	return &dispatcherHandlerInstance
}

//	配置控制器 (不是线程安全的)
//	projectPrefix	模块url  Prefix
//	ctrlPrefix		控制器url  ctrlPrefix
//	ctrl		控制器实例，可以是实例的指针或实例拷贝
func (dis *dispatcherHandler) PutCtrlPrefix(projectPrefix string,ctrlPrefix string,ctrl interface{}){
	if nil == dis.urlByProjectPrefixMapping[projectPrefix]{
		dis.urlByProjectPrefixMapping[projectPrefix] = make(map[string]interface{})
	}
	dis.urlByProjectPrefixMapping[projectPrefix][ctrlPrefix] = ctrl
}

//	为指定 projectPrefix 配置拦截器
//	每个 projectPrefix 只有一个拦截器
//	拦截器在调用控制器之前调用BeforeHandler()，在之后调用AfterHandler()
func (dis *dispatcherHandler) PutProjectInterceptor(projectPrefix string,handlerInterceptor HandlerInterceptorInterface){
	dis.handlerInterceptorMapping[projectPrefix] = handlerInterceptor
}

//	过滤器配置
//	实现 HttpFilterInterface 接口
//	过滤器只会在请求之初调用一次
func (dis *dispatcherHandler) SetHttpFilter(filter HttpFilterInterface){
	dis.filter = filter
}

//	出现错误时的失败响应
//	404 处理
//	500	处理
//	实现	FailureResponseInterface 接口，出现错误回调Failure404()、Failure500()
func (dis *dispatcherHandler) SetFailureResponse(failure FailureResponseInterface)  {
	dis.failure = failure
}

func (dis dispatcherHandler) ServeHTTP(writer http.ResponseWriter,request *http.Request){
	dis.HandlerFun()(writer,request)
}

//	请求转发HandleFun
//	调用顺序 ： 接收请求 ——> 过滤器 ->	拦截器BeforeHandler() -> 处理请求的控制器ctrl -> 拦截器AfterHandler()
//	请求路径寻找 ： projectPrefix/ctrlPrefix/方法名
func (dis *dispatcherHandler) HandlerFun() func (writer http.ResponseWriter, request *http.Request){
	return func (writer http.ResponseWriter, request *http.Request) {

		if nil != dis.filter{
			flag := dis.filter.Filter(writer,request)
			if !flag{
				return
			}
		}

		path := request.URL.Path
		urlSplit := strings.Split(path,"/")
		if 3 > len(urlSplit) {
			log.Println(fmt.Sprintf("path :[ %s ] 3 > len(urlSplit) unable to parse ",path))
			if nil != dis.failure{
				dis.failure.Failure500(writer,request)
			}else{
				writer.Write([]byte("Handler: ERROR URL FAIL！"))
			}
			return
		}
		if "" == urlSplit[0]{
			urlSplit = urlSplit[1:]
		}

		projectUrl := urlSplit[0]
		interceptor := dis.handlerInterceptorMapping[projectUrl]
		if nil != interceptor{
			flag,data := interceptor.BeforeHandler(writer,request)
			if !flag{
				writer.Write([]byte(data))
				return
			}
		}

		instanceMapping := dis.urlByProjectPrefixMapping[projectUrl]
		if nil == instanceMapping{
			log.Println(fmt.Sprintf("%s instanceMapping Not Found",projectUrl))
			failure404(dis.failure,writer,request)
			return
		}

		instanceUrl := urlSplit[1]
		if 3 < len(urlSplit){
			for i:=2;i<len(urlSplit);i++{
				instanceUrl = fmt.Sprintf("%s/%s",instanceUrl,urlSplit[i])
			}
		}
		instance := instanceMapping[instanceUrl]
		if nil == instance{
			log.Println(fmt.Sprintf("%s Instance Not Found",instanceUrl))
			failure404(dis.failure,writer,request)
			return
		}
		refValue := reflect.ValueOf(instance)
		methodName := strings.TrimSpace(urlSplit[len(urlSplit)-1])
		var refMethod reflect.Value
		refMethod = refValue.MethodByName(methodName)
		if !refMethod.IsValid() {
			refMethod = refValue.Elem().MethodByName(methodName)
		}

		if !refMethod.IsValid() {
			log.Println(fmt.Sprintf("%s Method Not Found",methodName))
			failure404(dis.failure,writer,request)
			return
		}

		refMtdType := refMethod.Type()
		numIn := refMtdType.NumIn()
		methodParams := make([]reflect.Value,numIn)
		if 1 == numIn{
			request.ParseForm()
			form := request.Form
			fiType := refMtdType.In(0)
			methodParams[0] = formToTypeValue(fiType,form)
		}else if 2 == numIn{
			methodParams[0] = reflect.ValueOf(writer)
			methodParams[1] = reflect.ValueOf(request)
		}else if 3 == numIn{
			request.ParseForm()
			form := request.Form
			fiType := refMtdType.In(0)
			methodParams[0] = formToTypeValue(fiType,form)
			methodParams[1] = reflect.ValueOf(writer)
			methodParams[2] = reflect.ValueOf(request)
		}

		result := refMethod.Call(methodParams)
		if nil != result && 0 < len(result){
			rsu := result[0]
			switch rsu.Kind() {
			case reflect.Struct,reflect.Map,reflect.Slice:
				marshalData,_ := json.Marshal(rsu.Interface())
				writer.Write(marshalData)
			default:
				writer.Write([]byte(rsu.String()))
			}
		}

		if nil != interceptor{
			interceptor.AfterHandler(writer,request)
		}
	}
}

func failure404(failure FailureResponseInterface,writer http.ResponseWriter, request *http.Request){
	if nil != failure{
		failure.Failure404(writer,request)
	}else{
		writer.Write([]byte("404"))
	}
}


//	map数据根据reflect.Type转为reflect.Value
//	fiType reflect.Type	类型,map、struct支持，其它都为string
//	form map[string][]string 数据源 如request.Form
//	reflect.Value	值
func formToTypeValue(fiType reflect.Type,form map[string][]string) reflect.Value{
	fiTypeKind := fiType.Kind()
	switch fiTypeKind {
	case reflect.Map:
		valMap := make(map[string]string,len(form))
		for key,values := range form {
			valMap[key] = stringArrayToString(values)
		}
		return reflect.ValueOf(valMap)
	case reflect.Struct:
		stVal := reflect.New(fiType)
		stType := stVal.Type()
		stElem := stType.Elem()
		numFiled := stElem.NumField()
		for i:=0;i<numFiled;i++{
			tf := stElem.Field(i)
			tagJson := tf.Tag.Get("json")
			for key,value := range form{
				if key == tagJson{
					v := stringToType(tf.Type.Name(),stringArrayToString(value))
					stVal.Elem().Field(i).Set(reflect.ValueOf(v))
					break
				}
			}
		}
		return stVal.Elem()
	default:
		vaList := ""
		for _,values := range form{
			vaList = stringArrayToString(values)
		}
		return reflect.ValueOf(vaList)
	}
}

//	将string参数转为typeStr指定类型的值
//	typeStr string	类型字串	支持int、float、bool、Time
//	valueStr string	值
//	interface{}	为 nil则失败
func stringToType(typeStr string,valueStr string) interface{}{
	var data interface{}
	var e error
	switch typeStr {
	case "int":
		data,e = strconv.Atoi(valueStr)
	case "int8":
		data,e = strconv.ParseInt(valueStr,10,8)
		if nil == e{
			data = data.(int8)
		}
	case "int16":
		data,e = strconv.ParseInt(valueStr,10,16)
		if nil == e{
			data = data.(int16)
		}
	case "int32":
		data,e = strconv.ParseInt(valueStr,10,32)
		if nil == e{
			data = data.(int32)
		}
	case "int64":
		data,e = strconv.ParseInt(valueStr,10,64)
	case "bool":
		data,e = strconv.ParseBool(valueStr)
	case "float32":
		data,e = strconv.ParseFloat(valueStr,32)
		if nil == e {
			data = float32(data.(float64))
		}
	case "float64":
		data,e = strconv.ParseFloat(valueStr,64)
	case "string":
		data = valueStr
	case "Time":
		if 10 == len(valueStr){
			data,e = time.Parse("2006-01-02",valueStr)
		}else if 13 == len(valueStr){
			data,e = time.Parse("2006-01-02 15",valueStr)
		}else if 16 == len(valueStr){
			data,e = time.Parse("2006-01-02 15:04",valueStr)
		}else if 19 == len(valueStr){
			data,e = time.Parse("2006-01-02 15:04:05",valueStr)
		}
		if nil != e {
			e = nil
			data = time.Now()
		}
	}
	if nil != e{
		log.Println(e)
		return nil
	}
	return data
}

//	字串数组转字串，以,拼接
//	strArr []string	字串数组
//	string	以,间隔的值
func stringArrayToString(strArr []string) string{
	str := ""
	for inx,_:= range strArr{
		if 0 == inx {
			str = strArr[inx]
			continue
		}
		str = fmt.Sprintf("%s,%s",str,strArr[inx])
	}
	return str
}
