package example

import (
	"fmt"
	"net/http"
	"text/template"
)

//  author: OD
//  since: 2023/7/15
//  desc:

type XSSFilterImplements struct {
}

// Intercept 只会过滤进行参数装载的控制器，如果是 (writer http.ResponseWriter, request *http.Request) 控制器，不会进行字符串过滤
// key, val string 字段的名称与原始值
// 返回 true 则拦截下请求不会继续分发，会就此中断
func (that *XSSFilterImplements) Intercept(key, val string, writer http.ResponseWriter, request *http.Request) bool {

	hs := template.HTMLEscapeString(val)
	if len(hs) != len(val) {
		fmt.Println(fmt.Sprintf("key=%s | val=%s | xssVal=%s", key, val, hs))
		return true
	}

	hs = template.JSEscapeString(val)
	if len(hs) != len(val) {
		fmt.Println(fmt.Sprintf("key=%s | val=%s | xssVal=%s", key, val, hs))
		return true
	}

	return false
}
