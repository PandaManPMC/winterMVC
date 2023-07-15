package winterMVC

import "net/http"

//  author: OD
//  since: 2023/7/15
//  desc:

type StringFilterInterface interface {
	// Intercept 只会过滤进行参数装载的控制器，如果是 (writer http.ResponseWriter, request *http.Request) 控制器，不会进行字符串过滤
	// key, val string 字段的名称与原始值
	// 返回 true 则拦截下请求不会继续分发，会就此中断
	Intercept(key, val string, writer http.ResponseWriter, request *http.Request) bool
}
