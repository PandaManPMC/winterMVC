package example

import (
	"fmt"
	winterMVC "github.com/PandaManPMC/winterMVC"
	"net/http"
	"testing"
	"time"
)

//	测试
//	访问方式	xxx:port/项目url/控制器映射键/方法名
//	http://localhost:7080/example/test/QueryList
//	http://localhost:7080/example/test/QueryListMap?name=heih&age=99&activity=true&fighting=33.55&inDate=2022-08-31T09:08:29.837820+00:00
//	http://localhost:7080/example/test/QueryListWR?name=heih&age=99&activity=true&fighting=33.55&inDate=2022-08-31T09:08:29.837820+00:00
//	http://localhost:7080/example/test/QueryListStruct?name=heih&age=99&activity=true&fighting=33.55&inDate=2022-08-31T09:08:29.837820+00:00
//	http://localhost:7080/example/test/QueryListStructWR?name=heih&age=99&activity=true&fighting=33.55&inDate=2022-08-31T09:08:29.837820+00:00

//	POST 测试 Content-Type application/json
//	http://localhost:7080/example/test/QueryListStruct
//	{"name":"laoniqiu","age":0,"activity":true,"fighting":33.55,"inDate":"2022-08-31T09:08:29.837820+00:00"}

//	POST 测试 Content-Type application/x-www-form-urlencoded
//	time 类型可以传递 2006-01-02 15:04:05 格式，无法解析的格式会报错

func TestMVC(t *testing.T) {

	//	获得mvc控制器 实例
	mvc := winterMVC.GetInstanceByDispatcherHandler()
	projectPrefix := "example"
	testC := GetInstanceByTestCtrl()
	//	存入控制器
	mvc.RouteCtrl(projectPrefix, "test", testC)
	//	http://localhost:7080/v2/example/test2/GetList
	mvc.RouteCtrl("v2/example", "test2", GetInstanceByTest2Ctrl())

	//	配置拦截器
	var inter myInterceptor
	mvc.RouteProjectInterceptor(projectPrefix, inter)

	//	配置过滤器
	var f filter
	mvc.SetHttpFilter(f)

	//	错误处理
	var fail failureResp
	mvc.SetFailureResponse(fail)

	//  配置日志输出
	var lo logs
	mvc.SetLogs(&lo)

	//	参数封装错误回调
	mvc.SetParameterError(&ParameterErrorImp{})

	// xss
	// 测试 http://localhost:7080/example/test/QueryListStruct?name=%3Cscript%3Ealert(1)%3C/script%3E&age=99&activity=true&fighting=33.55&inDate=2022-08-31T09:08:29.837820+00:00

	//	POST 测试 Content-Type application/json
	//	http://localhost:7080/example/test/QueryListStruct
	//	{"name":"<img src='http://google.com?token=abc'/>","age":0,"activity":true,"fighting":33.55,"inDate":"2022-08-31T09:08:29.837820+00:00"}
	mvc.SetStringFilterInterface(&XSSFilterImplements{})

	//	启动http服务 方式1
	//http.HandleFunc("/winterMvc/", mvc.HandlerFun())
	//http.Handle("/favicon.ico",http.FileServer(http.Dir("./web/img")))
	//http.ListenAndServe(":7080",nil)

	//	启动http服务 方式2
	maxHeaderBytes := 1024 * 1024 * 20
	server := http.Server{
		Addr: ":7080",
		//Handler: mvc,
		ReadTimeout:    time.Second * 60,
		WriteTimeout:   time.Second * 60,
		IdleTimeout:    time.Second * 60,
		MaxHeaderBytes: maxHeaderBytes,
	}
	http.Handle("/favicon.ico", http.FileServer(http.Dir("./web/img")))
	http.Handle("/", mvc)
	fmt.Println(server.ListenAndServe())
}
