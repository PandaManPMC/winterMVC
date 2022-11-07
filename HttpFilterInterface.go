package winterMVC

import "net/http"

//	http过滤器
//	Filter 过滤任何请求
type HttpFilterInterface interface {

	//	Filter 过滤器
	//	writer *http.ResponseWriter
	//	request *http.Request
	//	bool	true继续执行，false中断本次请求
	Filter(writer *http.ResponseWriter, request *http.Request) bool
}
