package winterMVC

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//	dispatcherHandler 核心控制器，接收请求分发处理
type dispatcherHandler struct {
	routeMapping              map[string]map[string]interface{}
	handlerInterceptorMapping map[string]HandlerInterceptorInterface
	filter                    HttpFilterInterface
	failure                   FailureResponseInterface
	logs                      LogsInterface
}

var dispatcherHandlerInstance dispatcherHandler

func init() {
	dispatcherHandlerInstance.routeMapping = make(map[string]map[string]interface{})
	dispatcherHandlerInstance.handlerInterceptorMapping = make(map[string]HandlerInterceptorInterface)
}

//	GetInstanceByDispatcherHandler 获取 MVC 实例
//	dispatcherHandler 实例指针
//	获得dispatcherHandler 初始化MVC
//	初始化方式1	http.HandleFunc("/projectRoute/", winter-mvc-core.HandlerFun())
//	初始化方式2	server := http.Server{ Handler: winter_mvc_core.GetInstanceByDispatcherHandler(), ...
func GetInstanceByDispatcherHandler() *dispatcherHandler {
	return &dispatcherHandlerInstance
}

//	SetLogs 配置日志输出
func (dis *dispatcherHandler) SetLogs(log LogsInterface) {
	dis.logs = log
}

func logInfo(msg string) {
	if nil != dispatcherHandlerInstance.logs {
		dispatcherHandlerInstance.logs.Info(msg)
	}
}

func logError(msg string, err interface{}) {
	if nil != dispatcherHandlerInstance.logs {
		dispatcherHandlerInstance.logs.Error(msg, err)
	}
}

//	RouteCtrl 配置控制器 (不是线程安全的)
//	projectRoute	模块 url  Prefix
//	ctrlRoute		控制器 url  ctrlPrefix
//	ctrlInstance	控制器实例，可以是实例的指针或实例拷贝
func (dis *dispatcherHandler) RouteCtrl(projectRoute string, ctrlRoute string, ctrlInstance interface{}) {
	if nil == dis.routeMapping[projectRoute] {
		dis.routeMapping[projectRoute] = make(map[string]interface{})
	}
	dis.routeMapping[projectRoute][ctrlRoute] = ctrlInstance
}

//	RouteProjectInterceptor 为指定 projectRoute 配置拦截器
//	每个 projectPrefix 只有一个拦截器
//	拦截器在调用控制器之前调用BeforeHandler()，在之后调用AfterHandler()
func (dis *dispatcherHandler) RouteProjectInterceptor(projectRoute string, handlerInterceptor HandlerInterceptorInterface) {
	dis.handlerInterceptorMapping[projectRoute] = handlerInterceptor
}

//	SetHttpFilter 过滤器配置
//	实现 HttpFilterInterface 接口
//	过滤器只会在请求之初调用一次
func (dis *dispatcherHandler) SetHttpFilter(filter HttpFilterInterface) {
	dis.filter = filter
}

//	SetFailureResponse 出现错误时的失败响应
//	404 处理
//	500	处理
//	实现	FailureResponseInterface 接口，出现错误回调Failure404()、Failure500()
func (dis *dispatcherHandler) SetFailureResponse(failure FailureResponseInterface) {
	dis.failure = failure
}

func (dis dispatcherHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	dis.HandlerFun()(writer, request)
}

//	deferRecover 异常捕获
func deferRecover(dis *dispatcherHandler, writer http.ResponseWriter, request *http.Request) {
	err := recover()
	if nil == err {
		return
	}
	logError("deferRecover", err)
	if nil != dis.failure {
		dis.failure.Failure500(writer, request)
	} else {
		writer.Write([]byte("SERVER: ERROR！"))
	}
}

//	HandlerFun 请求转发
//	调用顺序 ： 接收请求 ——> 过滤器 ->	拦截器BeforeHandler() -> 处理请求的控制器ctrl -> 拦截器AfterHandler()
//	请求路径寻找 ： uri/projectRoute/ctrlRoute/ctrlInstance方法名
func (dis *dispatcherHandler) HandlerFun() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer deferRecover(dis, writer, request)
		if nil != dis.filter {
			flag := dis.filter.Filter(writer, request)
			if !flag {
				return
			}
		}

		path := request.URL.Path
		urlSplit := strings.Split(path, "/")
		if 3 > len(urlSplit) {
			logInfo(fmt.Sprintf("path :[ %s ] 3 > len(urlSplit) unable to parse ", path))
			if nil != dis.failure {
				dis.failure.Failure500(writer, request)
			} else {
				writer.Write([]byte("Handler: ERROR URL FAIL！"))
			}
			return
		}
		if "" == urlSplit[0] {
			urlSplit = urlSplit[1:]
		}

		projectUrl := urlSplit[0]
		interceptor := dis.handlerInterceptorMapping[projectUrl]
		if nil != interceptor {
			flag, data := interceptor.BeforeHandler(writer, request)
			if !flag && "" != data {
				writer.Write([]byte(data))
				return
			}
			if !flag {
				return
			}
		}

		instanceMapping := dis.routeMapping[projectUrl]
		if nil == instanceMapping {
			logInfo(fmt.Sprintf("%s instanceMapping Not Found", projectUrl))
			failure404(dis.failure, writer, request)
			return
		}

		instanceUrl := urlSplit[1]
		if 3 < len(urlSplit) {
			for i := 2; i < len(urlSplit); i++ {
				instanceUrl = fmt.Sprintf("%s/%s", instanceUrl, urlSplit[i])
			}
		}
		instance := instanceMapping[instanceUrl]
		if nil == instance {
			logInfo(fmt.Sprintf("%s Instance Not Found", instanceUrl))
			failure404(dis.failure, writer, request)
			return
		}
		refValue := reflect.ValueOf(instance)
		methodName := strings.TrimSpace(urlSplit[len(urlSplit)-1])
		var refMethod reflect.Value
		refMethod = refValue.MethodByName(methodName)
		if !refMethod.IsValid() && reflect.Ptr == refValue.Kind() {
			refMethod = refValue.Elem().MethodByName(methodName)
		}

		if !refMethod.IsValid() {
			logInfo(fmt.Sprintf("%s Method Not Found", methodName))
			failure404(dis.failure, writer, request)
			return
		}

		refMtdType := refMethod.Type()
		numIn := refMtdType.NumIn()
		methodParams := make([]reflect.Value, numIn)

		if 0 != numIn {
			requestParams(writer, request, &methodParams, refMtdType)
		}
		// 响应
		result := refMethod.Call(methodParams)
		if nil != result && 0 < len(result) {
			rsu := result[0]
			switch rsu.Kind() {
			case reflect.Struct, reflect.Map, reflect.Slice:
				marshalData, _ := json.Marshal(rsu.Interface())
				writer.Write(marshalData)
			default:
				writer.Write([]byte(rsu.String()))
			}
		}

		if nil != interceptor {
			interceptor.AfterHandler(writer, request)
		}
	}
}

//	requestParams 请求参数封装
func requestParams(writer http.ResponseWriter, request *http.Request, methodParams *[]reflect.Value, refMtdType reflect.Type) {
	for i := 0; i < refMtdType.NumIn(); i++ {
		inType := refMtdType.In(i)
		switch inType.String() {
		case "*http.Request":
			(*methodParams)[i] = reflect.ValueOf(request)
		case "http.ResponseWriter":
			(*methodParams)[i] = reflect.ValueOf(writer)
		default:
			requestToData(request, methodParams, inType, i)
		}
	}
}

// requestToData request 中参数封装进结构体或 map，支持 Content-Type 【application/json || application/x-www-form-urlencoded】
func requestToData(request *http.Request, methodParams *[]reflect.Value, inType reflect.Type, index int) {
	if ContentTypeIsJSON(request) {
		// 以 【application/json】
		defer request.Body.Close()
		buf, err := ioutil.ReadAll(request.Body)
		if nil != err {
			logError("requestToData", err)
			return
		}
		obj := reflect.New(inType)
		json.Unmarshal(buf, obj.Interface())
		(*methodParams)[index] = obj.Elem()
	} else {
		// 其它-以 form 形式读取参数 【application/x-www-form-urlencoded】
		request.ParseForm()
		form := request.Form
		(*methodParams)[index] = formToTypeValue(inType, form)
	}
}

func ContentTypeIsJSON(request *http.Request) bool {
	contentType := GetContentType(request)
	if "application/json" == contentType {
		return true
	}
	return false
}

func GetContentType(request *http.Request) string {
	contentType := request.Header.Get("Content-Type")
	if 0 == len(contentType) {
		return request.Header.Get("content-type")
	}
	return contentType
}

func failure404(failure FailureResponseInterface, writer http.ResponseWriter, request *http.Request) {
	if nil != failure {
		failure.Failure404(writer, request)
	} else {
		writer.Write([]byte("404"))
	}
}

//	map数据根据reflect.Type转为reflect.Value
//	fiType reflect.Type	类型,map、struct支持，其它都为string
//	form map[string][]string 数据源 如request.Form
//	reflect.Value	值
func formToTypeValue(fiType reflect.Type, form map[string][]string) reflect.Value {
	fiTypeKind := fiType.Kind()
	switch fiTypeKind {
	case reflect.Map:
		valMap := make(map[string]string, len(form))
		for key, values := range form {
			valMap[key] = stringArrayToString(values)
		}
		return reflect.ValueOf(valMap)
	case reflect.Struct:
		stVal := reflect.New(fiType)
		stType := stVal.Type()
		stElem := stType.Elem()
		numFiled := stElem.NumField()
		for i := 0; i < numFiled; i++ {
			tf := stElem.Field(i)
			tagJson := tf.Tag.Get("json")
			for key, value := range form {
				if key == tagJson {
					if 0 == len(value) {
						break
					}
					sv := stringArrayToString(value)
					if "" == sv && tf.Type.Name() != "string" {
						break
					}
					v := stringToType(tf.Type.Name(), sv)
					stVal.Elem().Field(i).Set(reflect.ValueOf(v))
					break
				}
			}
		}
		return stVal.Elem()
	default:
		vaList := ""
		for _, values := range form {
			vaList = stringArrayToString(values)
		}
		return reflect.ValueOf(vaList)
	}
}

//	将string参数转为typeStr指定类型的值
//	typeStr string	类型字串	支持int、float、bool、Time
//	valueStr string	值
//	interface{}	为 nil则失败
func stringToType(typeStr string, valueStr string) interface{} {
	var data interface{}
	var e error
	switch typeStr {
	case "int":
		if "" == valueStr {
			return 0
		}
		data, e = strconv.Atoi(valueStr)
	case "uint":
		if "" == valueStr {
			return uint(0)
		}
		val, e1 := strconv.Atoi(valueStr)
		if nil == e1 {
			data = uint(val)
		} else {
			e = e1
			data = val
		}
	case "int8":
		if "" == valueStr {
			return int8(0)
		}
		data, e = strconv.ParseInt(valueStr, 10, 8)
		if nil == e {
			data = int8(data.(int64))
		}
	case "uint8":
		if "" == valueStr {
			return uint8(0)
		}
		data, e = strconv.ParseUint(valueStr, 10, 8)
		if nil == e {
			data = uint8(data.(uint64))
		}
	case "int16":
		if "" == valueStr {
			return int16(0)
		}
		data, e = strconv.ParseInt(valueStr, 10, 16)
		if nil == e {
			data = int16(data.(int64))
		}
	case "uint16":
		if "" == valueStr {
			return uint16(0)
		}
		data, e = strconv.ParseUint(valueStr, 10, 16)
		if nil == e {
			data = uint16(data.(uint64))
		}
	case "int32":
		if "" == valueStr {
			return int32(0)
		}
		data, e = strconv.ParseInt(valueStr, 10, 32)
		if nil == e {
			data = int32(data.(int64))
		}
	case "uint32":
		if "" == valueStr {
			return uint32(0)
		}
		data, e = strconv.ParseUint(valueStr, 10, 32)
		if nil == e {
			data = uint32(data.(uint64))
		}
	case "int64":
		if "" == valueStr {
			return int64(0)
		}
		data, e = strconv.ParseInt(valueStr, 10, 64)
	case "uint64":
		if "" == valueStr {
			return uint64(0)
		}
		data, e = strconv.ParseUint(valueStr, 10, 64)
	case "bool":
		if "" == valueStr {
			return false
		}
		data, e = strconv.ParseBool(valueStr)
	case "float32":
		if "" == valueStr {
			return 0
		}
		data, e = strconv.ParseFloat(valueStr, 32)
		if nil == e {
			data = float32(data.(float64))
		}
	case "float64":
		if "" == valueStr {
			return 0
		}
		data, e = strconv.ParseFloat(valueStr, 64)
	case "string":
		if "" == valueStr {
			return ""
		}
		data = valueStr
	case "Time":
		if 10 == len(valueStr) {
			data, e = time.Parse("2006-01-02", valueStr)
		} else if 13 == len(valueStr) {
			data, e = time.Parse("2006-01-02 15", valueStr)
		} else if 16 == len(valueStr) {
			data, e = time.Parse("2006-01-02 15:04", valueStr)
		} else if 19 == len(valueStr) {
			data, e = time.Parse("2006-01-02 15:04:05", valueStr)
		}
		if nil != e {
			e = nil
			data = time.Now()
		}
	}
	if nil != e {
		logError("stringToType", e)
		return nil
	}
	return data
}

//	字串数组转字串，以,拼接
//	strArr []string	字串数组
//	string	以【】间隔的值
func stringArrayToString(strArr []string) string {
	str := ""
	for inx, _ := range strArr {
		if 0 == inx {
			str = strArr[inx]
			continue
		}
		str = fmt.Sprintf("%s【】%s", str, strArr[inx])
	}
	return str
}
