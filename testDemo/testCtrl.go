package testDemo

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type testCtrl struct {
}
var testCtrlInstance testCtrl
func GetInstanceByTestCtrl() *testCtrl{
	return &testCtrlInstance
}

func (t *testCtrl) QueryList() result{
	fmt.Println("调用了 QueryList 方法")
	r := resultNewSuccess("成功 QueryList",111111111111111)
	fmt.Println("QueryList 响应：",r)
	return r
}

func (t *testCtrl) QueryList2(params map[string]string) result{
	fmt.Println("调用了 QueryList2 方法")

	idc := params["IdCategory"]
	fmt.Println("idc=",idc)
	if "" != idc {
		iint,_ := strconv.ParseInt(idc,10,64)
		fmt.Println("成功得到值=",iint)
	}

	r := resultNewSuccess("成功 QueryList2",params)
	fmt.Println("QueryList 响应：",r)
	return r
}

func (t *testCtrl) QueryList3(w http.ResponseWriter, request *http.Request) result{
	fmt.Println("调用了 QueryList3 方法")
	fo := request.Form
	fmt.Println(fo)
	r := resultNewSuccess("成功  QueryList3",fo)
	fmt.Println("QueryList 响应：",r)
	return r
}

type dog struct {
	Name string	`json:"name" table:"name_name"`
	Age int	`json:"age"`
	Fighting float32 `json:"fighting"`
	Activity bool `json:"activity"`
	InDate time.Time `json:"inDate"`
}

func (t *testCtrl) QueryList4(dg dog) result{
	fmt.Println("调用了 QueryList4 方法")
	fmt.Println(dg)
	r := resultNewSuccess("成功  QueryList4",dg)
	fmt.Println("QueryList4 响应：",r)
	return r
}

func (t *testCtrl) QueryList5(dg dog,w http.ResponseWriter, request *http.Request) result{
	fmt.Println("调用了 QueryList5 方法")
	fmt.Println(dg)
	head := request.Header
	fmt.Println(head)
	w.Header().Set("testHeader","哈哈哈哈哈 响应头")
	r := resultNewSuccess("成功  QueryList5",dg)
	fmt.Println("QueryList5 响应：",r)
	return r
}
