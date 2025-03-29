package winterMVC

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// dispatcherHandler 核心控制器，接收请求分发处理
type dispatcherHandler struct {
	routeMapping              map[string]map[string]interface{}
	handlerInterceptorMapping map[string]HandlerInterceptorInterface
	filter                    HttpFilterInterface
	failure                   FailureResponseInterface
	logs                      LogsInterface
	parameterError            ParameterErrorInterface
	stringFilter              StringFilterInterface // 字符串参数过滤回调，如果没有设置则默认不拦截
}

var dispatcherHandlerInstance dispatcherHandler

func init() {
	dispatcherHandlerInstance.routeMapping = make(map[string]map[string]interface{})
	dispatcherHandlerInstance.handlerInterceptorMapping = make(map[string]HandlerInterceptorInterface)
}

// GetInstanceByDispatcherHandler 获取 MVC 实例
// dispatcherHandler 实例指针
// 获得dispatcherHandler 初始化MVC
// 初始化方式1	http.HandleFunc("/projectRoute/", winter-mvc-core.HandlerFun())
// 初始化方式2	server := http.Server{ Handler: winter_mvc_core.GetInstanceByDispatcherHandler(), ...
func GetInstanceByDispatcherHandler() *dispatcherHandler {
	return &dispatcherHandlerInstance
}

// SetLogs 配置日志输出
func (dis *dispatcherHandler) SetLogs(log LogsInterface) {
	dis.logs = log
}

// SetStringFilterInterface 设置 字符串 拦截接口
func (dis *dispatcherHandler) SetStringFilterInterface(xf StringFilterInterface) {
	dis.stringFilter = xf
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

// RouteCtrl 配置控制器 (不是线程安全的)
// projectRoute	模块 url  Prefix
// ctrlRoute		控制器 url  ctrlPrefix
// ctrlInstance	控制器实例，可以是实例的指针或实例拷贝
func (dis *dispatcherHandler) RouteCtrl(projectRoute string, ctrlRoute string, ctrlInstance interface{}) {
	if nil == dis.routeMapping[projectRoute] {
		dis.routeMapping[projectRoute] = make(map[string]interface{})
	}
	dis.routeMapping[projectRoute][ctrlRoute] = ctrlInstance
}

// RouteProjectInterceptor 为指定 projectRoute 配置拦截器
// 每个 projectPrefix 只有一个拦截器
// 拦截器在调用控制器之前调用BeforeHandler()，在之后调用AfterHandler()
func (dis *dispatcherHandler) RouteProjectInterceptor(projectRoute string, handlerInterceptor HandlerInterceptorInterface) {
	dis.handlerInterceptorMapping[projectRoute] = handlerInterceptor
}

// SetHttpFilter 过滤器配置
// 实现 HttpFilterInterface 接口
// 过滤器只会在请求之初调用一次
func (dis *dispatcherHandler) SetHttpFilter(filter HttpFilterInterface) {
	dis.filter = filter
}

// SetFailureResponse 出现错误时的失败响应
// 404 处理
// 500	处理
// 实现	FailureResponseInterface 接口，出现错误回调Failure404()、Failure500()
func (dis *dispatcherHandler) SetFailureResponse(failure FailureResponseInterface) {
	dis.failure = failure
}

// SetParameterError 参数装载发生错误的回调
// callback ParameterErrorInterface 回调，实现 ParameterErrorInterface
func (dis *dispatcherHandler) SetParameterError(callback ParameterErrorInterface) {
	dis.parameterError = callback
}

func (dis *dispatcherHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	dis.HandlerFun()(writer, request)
}

// deferRecover 异常捕获
func deferRecover(dis *dispatcherHandler, writer http.ResponseWriter, request *http.Request) {
	err := recover()
	if nil == err {
		return
	}
	if nil != dis.failure {
		dis.failure.Failure500(writer, request, err)
	} else {
		logError("deferRecover", err)
		writer.Write([]byte("SERVER: ERROR！"))
	}
}

// HandlerFun 请求转发
// 调用顺序 ： 接收请求 ——> 过滤器 ->	拦截器BeforeHandler() -> 处理请求的控制器ctrl -> 拦截器AfterHandler()
// 请求路径寻找 ： uri/projectRoute/ctrlRoute/ctrlInstance方法名
func (dis *dispatcherHandler) HandlerFun() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer deferRecover(dis, writer, request)
		if nil != dis.filter {
			flag := dis.filter.Filter(&writer, request)
			if !flag {
				return
			}
		}

		path := request.URL.Path
		urlSplit := strings.Split(path, "/")
		if 3 > len(urlSplit) {
			logInfo(fmt.Sprintf("path :[ %s ] 3 > len(urlSplit) unable to parse ", path))
			if nil != dis.failure {
				dis.failure.Failure500(writer, request, errors.New("route error"))
			} else {
				writer.Write([]byte("Handler: ERROR URL FAIL！"))
			}
			return
		}
		if "" == urlSplit[0] {
			urlSplit = urlSplit[1:]
		}

		projectUrl := urlSplit[0]
		// v 开头之后数字，则合并 urlSplit[0]/urlSplit[1]
		if MustCompileVersion(projectUrl) {
			if 4 > len(urlSplit) {
				logInfo(fmt.Sprintf("path :[ %s ] 4 > len(urlSplit) unable to parse ", path))
				if nil != dis.failure {
					dis.failure.Failure500(writer, request, errors.New("route error"))
				} else {
					writer.Write([]byte("Handler: ERROR URL FAIL！"))
				}
				return
			}
			projectUrl = fmt.Sprintf("%s/%s", urlSplit[0], urlSplit[1])
			urlSplit = urlSplit[1:]
		}

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
			if err := dis.requestParams(writer, request, &methodParams, refMtdType, path); nil != err {
				logError(err.Error(), nil)
				if nil != dis.parameterError {
					dis.parameterError.ParameterError(writer, request, err)
				}
				return
			}
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

// requestParams 请求参数封装
func (dis *dispatcherHandler) requestParams(writer http.ResponseWriter, request *http.Request, methodParams *[]reflect.Value, refMtdType reflect.Type, path string) error {
	for i := 0; i < refMtdType.NumIn(); i++ {
		inType := refMtdType.In(i)
		switch inType.String() {
		case "*http.Request":
			(*methodParams)[i] = reflect.ValueOf(request)
		case "http.ResponseWriter":
			(*methodParams)[i] = reflect.ValueOf(writer)
		default:
			if err := dis.requestToData(writer, request, methodParams, inType, i, path); nil != err {
				return err
			}
		}
	}
	return nil
}

// requestToData request 中参数封装进结构体或 map，支持 Content-Type 【application/json || application/x-www-form-urlencoded】
func (dis *dispatcherHandler) requestToData(writer http.ResponseWriter, request *http.Request, methodParams *[]reflect.Value, inType reflect.Type, index int, path string) error {
	if ContentTypeIsJSON(request) {
		// 以 【application/json】
		defer request.Body.Close()
		buf, err := io.ReadAll(request.Body)
		if nil != err {
			logError(fmt.Sprintf("requestToData=%s", string(buf)), err)
			return errors.New("request to data failure")
		}
		obj := reflect.New(inType)
		if err := json.Unmarshal(buf, obj.Interface()); nil != err {
			logError(fmt.Sprintf("requestToData=%s", string(buf)), err)
			return errors.New("request to json Unmarshal data failure")
		}
		(*methodParams)[index] = obj.Elem()

		if reflect.Slice == inType.Kind() {
			return nil
		}
		mp := make(map[string]interface{})
		json.Unmarshal(buf, &mp)
		return dis.requiredJSON(writer, request, obj, mp)
	} else {
		// 其它-以 form 形式读取参数 【application/x-www-form-urlencoded】
		request.ParseForm()
		form := request.Form
		var err error
		if (*methodParams)[index], err = dis.formToTypeValue(writer, request, inType, form); nil != err {
			return err
		}
	}
	return nil
}

// InterceptString 字符串参数拦截，当 error 返回不为 nil 时，打断这个分发
func (dis *dispatcherHandler) InterceptString(writer http.ResponseWriter, request *http.Request, value, key string) error {
	if nil == dis.stringFilter {
		return nil
	}
	if dis.stringFilter.Intercept(key, value, writer, request) {
		return errors.New("parameter invalid")
	}
	return nil
}

func (dis *dispatcherHandler) requiredJSON(writer http.ResponseWriter, request *http.Request, bean reflect.Value, mp map[string]interface{}) error {
	beanElem := bean.Elem()
	num := beanElem.NumField()
	t := beanElem.Type()
	for i := 0; i < num; i++ {
		f := t.Field(i)
		rd := f.Tag.Get("required")
		json := f.Tag.Get("json")
		val, isExist := mp[json]

		//	字符串默认去除首位空格
		if "string" == f.Type.Name() && isExist {
			trimSpace := f.Tag.Get("trimSpace")
			if "false" != trimSpace {
				field := beanElem.Field(i)
				field.Set(reflect.ValueOf(strings.TrimSpace(field.String())))
				val = strings.TrimSpace(val.(string))
			}
		}

		// xss 过滤
		if "string" == f.Type.Name() && isExist {
			if e := dis.InterceptString(writer, request, val.(string), json); nil != e {
				return e
			}
		}

		if "true" == rd {
			//field := beanElem.Field(i)
			if !isExist {
				return errors.New(fmt.Sprintf("missing required parameters by 【%s】", json))
			}
			mi := f.Tag.Get("min")
			if "" != mi {
				// 最小值，作用于字符串长度
				if "string" == f.Type.Name() {
					s := val.(string)
					minLen, _ := strconv.Atoi(mi)
					if minLen > StringToCharacterLen(s) {
						return errors.New(fmt.Sprintf("Parameter 【%s】 minimum length of %d, your %d", json, minLen, StringToCharacterLen(s)))
					}
				}
			}
		}
		if !isExist {
			// 非必传，如果值不存在，直接回退
			continue
		}
		// 最大值，作用于字符串即字符最大长度
		ma := f.Tag.Get("max")
		if "" != ma {
			if "string" == f.Type.Name() {
				s := val.(string)
				maxLen, _ := strconv.Atoi(ma)
				if maxLen < StringToCharacterLen(s) {
					return errors.New(fmt.Sprintf("parameter 【%s】 maximum length %d, yours %d", json, maxLen, StringToCharacterLen(s)))
				}
			}
		}
	}
	return nil
}

// StringToCharacterLen 获取字符串准确字符长度，而非字节长度
func StringToCharacterLen(s string) int {
	return len([]rune(s))
}

func ContentTypeIsJSON(request *http.Request) bool {
	contentType := GetContentType(request)
	if strings.Contains(contentType, "application/json") {
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

// formToTypeValue
// map数据根据reflect.Type转为reflect.Value
// fiType reflect.Type	类型,map、struct支持，其它都为string
// form map[string][]string 数据源 如request.Form
// reflect.Value	值
func (dis *dispatcherHandler) formToTypeValue(writer http.ResponseWriter, request *http.Request, fiType reflect.Type, form map[string][]string) (reflect.Value, error) {
	fiTypeKind := fiType.Kind()
	switch fiTypeKind {
	case reflect.Map:
		valMap := make(map[string]string, len(form))
		for key, values := range form {
			valMap[key] = stringArrayToString(values)

			// xss 过滤
			if e := dis.InterceptString(writer, request, valMap[key], key); nil != e {
				return reflect.ValueOf(valMap), e
			}
		}
		return reflect.ValueOf(valMap), nil
	case reflect.Struct:
		stVal := reflect.New(fiType)
		stType := stVal.Type()
		stElem := stType.Elem()
		numFiled := stElem.NumField()
		for i := 0; i < numFiled; i++ {
			tf := stElem.Field(i)
			tagJson := tf.Tag.Get("json")
			required := tf.Tag.Get("required")

			value, isExist := form[tagJson]
			if "true" == required {
				if !isExist {
					return stVal, errors.New(fmt.Sprintf("missing required parameters by 【%s】", tagJson))
				}
			} else {
				if !isExist {
					continue
				}
			}

			sv := stringArrayToString(value)
			if "" == sv && "string" != tf.Type.Name() {
				continue
			}

			if "string" == tf.Type.Name() {
				trimSpace := tf.Tag.Get("trimSpace")
				if "false" != trimSpace {
					sv = strings.TrimSpace(sv)
				}
			}

			// xss 过滤
			if "string" == tf.Type.Name() {
				if e := dis.InterceptString(writer, request, sv, tagJson); nil != e {
					return stVal, e
				}
			}

			// 字符串类型校验长度
			if "true" == required && "string" == tf.Type.Name() {
				mi := tf.Tag.Get("min")
				if "" != mi {
					minLen, _ := strconv.Atoi(mi)
					if minLen > StringToCharacterLen(sv) {
						return stVal, errors.New(fmt.Sprintf("Parameter 【%s】 minimum length of %d, your %d", tagJson, minLen, StringToCharacterLen(sv)))
					}
				}
			}

			ma := tf.Tag.Get("max")
			if "" != ma {
				if "string" == tf.Type.Name() {
					maxLen, _ := strconv.Atoi(ma)
					if maxLen < StringToCharacterLen(sv) {
						return stVal, errors.New(fmt.Sprintf("parameter 【%s】 maximum length %d, yours %d", tagJson, maxLen, StringToCharacterLen(sv)))
					}
				}
			}

			v := stringToType(tf.Type.Name(), sv)
			stVal.Elem().Field(i).Set(reflect.ValueOf(v))

			//for key, value := range form {
			//	if key == tagJson {
			//		if 0 == len(value) {
			//			if "true" == required {
			//				return stVal, errors.New(fmt.Sprintf("missing required parameters %s", tagJson))
			//			}
			//			break
			//		}
			//		sv := stringArrayToString(value)
			//		if "" == sv && tf.Type.Name() != "string" {
			//			break
			//		}
			//		v := stringToType(tf.Type.Name(), sv)
			//		stVal.Elem().Field(i).Set(reflect.ValueOf(v))
			//		break
			//	}
			//}
		}
		return stVal.Elem(), nil
	default:
		vaList := ""
		for key, values := range form {
			vaList = stringArrayToString(values)
			// xss 过滤
			if e := dis.InterceptString(writer, request, vaList, key); nil != e {
				return reflect.ValueOf(vaList), e
			}
		}
		return reflect.ValueOf(vaList), nil
	}
}
