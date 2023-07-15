package winterMVC

import "net/http"

//  author: OD
//  since: 2023/7/15
//  desc:

type XSSFilterInterface interface {
	// InterceptXSS 只会过滤进行参数装载的控制器，如果是 (writer http.ResponseWriter, request *http.Request) 控制器，不会进行 XSS 过滤
	// key, val , xssVal string xss 字段的名称与原始值、 xss 值
	// 返回 true 则拦截下请求不会继续分发，会就此中断
	InterceptXSS(key, val, xssVal string, writer http.ResponseWriter, request *http.Request) bool
}
