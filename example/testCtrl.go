package example

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type testCtrl struct {
}

var testCtrlInstance testCtrl

func GetInstanceByTestCtrl() *testCtrl {
	return &testCtrlInstance
}

func (t *testCtrl) QueryList() result {
	fmt.Println("调用了 QueryList 方法")
	r := resultNewSuccess("成功 QueryList", 111111111111111)
	fmt.Println("QueryList 响应：", r)
	return r
}

func (t *testCtrl) QueryListMap(params map[string]string) result {
	fmt.Println("调用了 QueryListMap 方法")

	idc := params["IdCategory"]
	fmt.Println("idc=", idc)
	if "" != idc {
		iint, _ := strconv.ParseInt(idc, 10, 64)
		fmt.Println("成功得到值=", iint)
	}

	r := resultNewSuccess("成功 QueryListMap", params)
	fmt.Println("QueryListMap 响应：", r)
	return r
}

func (t *testCtrl) QueryListWR(w http.ResponseWriter, request *http.Request) result {
	fmt.Println("调用了 QueryListWR 方法")
	fo := request.Form
	fmt.Println(fo)
	r := resultNewSuccess("成功  QueryListWR", fo)
	fmt.Println("QueryListWR 响应：", r)
	return r
}

type dog struct {
	Name     string    `json:"name" table:"name_name" required:"true"`
	Age      int       `json:"age" required:"true"`
	Fighting float32   `json:"fighting"`
	Activity bool      `json:"activity"`
	InDate   time.Time `json:"inDate"` // 只支持 rfc3339 格式，如 【2022-08-31T09:08:29.837820+00:00】
}

func (t *testCtrl) QueryListStruct(dg dog) result {
	fmt.Println("调用了 QueryListStruct 方法")
	fmt.Println(dg)
	r := resultNewSuccess("成功  QueryListStruct", dg)
	fmt.Println("QueryListStruct 响应：", r)
	return r
}

func (t *testCtrl) QueryListStructWR(dg dog, w http.ResponseWriter, request *http.Request) result {
	fmt.Println("调用了 QueryListStructWR 方法")
	fmt.Println(dg)
	head := request.Header
	fmt.Println(head)
	w.Header().Set("testHeader", "test response header")
	r := resultNewSuccess("成功  QueryListStructWR", dg)
	fmt.Println("QueryListStructWR 响应：", r)
	return r
}

//	http://127.0.0.1:7080/example/test/TestJSON
func (t *testCtrl) TestJSON(dg dog, w http.ResponseWriter) result {
	fmt.Println("调用了 TestJSON 方法")
	fmt.Println(dg)
	fmt.Println(w)
	return resultNewSuccess("成功  QueryListStructWR", dg)
}

type ami struct {
	Category string `json:"category"`
	Dog      dog    `json:"dog"`
}

// TestJSON2 测试嵌套结构体
func (t *testCtrl) TestJSON2(am ami, w http.ResponseWriter) result {
	//{
	//	"category": "哺乳",
	//	"dog": {
	//		"name": "黑子",
	//			"age": 88,
	//			"fighting": 33.22,
	//			"activity": true,
	//			"inDate": "2022-08-31T09:08:29.837820+00:00"
	//	}
	//}

	fmt.Println("调用了 TestJSON2 方法")
	fmt.Println(am)
	fmt.Println(w)
	return resultNewSuccess("成功  QueryListStructWR", am)
}
