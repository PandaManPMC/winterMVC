package testDemo

import (
	"fmt"
	"net/http"
)

//	错误处理器
//	实现 FailureResponseInterface 接口
type failureResp struct {

}


func (f failureResp) Failure404(writer http.ResponseWriter,request *http.Request){
	fmt.Println("出现 了Failure404 错误")
}

func (f failureResp) Failure500(writer http.ResponseWriter,request *http.Request){
	fmt.Println("出现 了Failure500 错误")
}
