package example

import (
	"fmt"
	"net/http"
)

//	过滤器
//	实现 HttpFilterInterface 接口
type filter struct {
}

func (f filter) Filter(writer *http.ResponseWriter, request *http.Request) bool {
	fmt.Println("过滤器执行。。。。。。。。。")
	return true
}
