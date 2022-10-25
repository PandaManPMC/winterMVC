package example

import (
	"fmt"
	"net/http"
)

//  author: laoniqiu
//  since: 2022/10/25
//  desc: example

type ParameterErrorImp struct {
}

func (*ParameterErrorImp) ParameterError(writer http.ResponseWriter, request *http.Request, err error) {
	fmt.Println("参数封装错误", err.Error())
	writer.Write([]byte(err.Error()))
}
