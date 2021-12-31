package winterMvc

import "net/http"

type FailureResponseInterface interface {

	//	错误处理回调404
	//	找不到对应控制器或处理方法回调404
	Failure404(writer http.ResponseWriter,request *http.Request)

	//	错误处理回调500
	//	出现预料之外的事，如url截取后长度不足3位
	Failure500(writer http.ResponseWriter,request *http.Request)

}


