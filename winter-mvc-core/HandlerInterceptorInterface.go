package winterMvc

import "net/http"

//	拦截器
//	实现接口 BeforeHandler()和 AfterHandler()
type HandlerInterceptorInterface interface {

	//	在处理器处理请求之前
	//	writer http.ResponseWriter
	//	request *http.Request
	//	bool	响应 true 执行处理，false不继续执行处理也不会执行AfterHandler()
	//	string	在bool为false拦截下请求时的响应数据
	BeforeHandler(writer http.ResponseWriter, request *http.Request) (bool,string)

	//	处理器处理请求之后
	//	writer http.ResponseWriter
	//	request *http.Request
	AfterHandler(writer http.ResponseWriter, request *http.Request)


}




