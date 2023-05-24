package example

import (
	"fmt"
	"net/http"
)

type test2Ctrl struct {
}

var test2CtrlInstance test2Ctrl

func GetInstanceByTest2Ctrl() *test2Ctrl {
	return &test2CtrlInstance
}

func (t *test2Ctrl) GetList(request *http.Request) result {
	fmt.Println("调用了 test2Ctrl GetList 方法")
	r := resultNewSuccess("成功 test2Ctrl GetList", 111111111111111)
	fmt.Println("test2Ctrl QueryList 响应：", r)
	return r
}
