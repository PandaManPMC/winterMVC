package example

import (
	"fmt"
	"net/http"
)

//	拦截器
//	实现 HandlerInterceptorInterface 接口
type myInterceptor struct {

}

func (m myInterceptor) BeforeHandler(writer http.ResponseWriter, request *http.Request) (bool,string){
	fmt.Println("BeforeHandler 执行.................................................")
	return true,"失败"
}

func (m myInterceptor) AfterHandler(writer http.ResponseWriter, request *http.Request){
	fmt.Println("AfterHandler 执行................................................")
}