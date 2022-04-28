package example

import (
	"fmt"
	"net/http"
	"testing"
	"time"
	"winterMVC"
)

//	测试
//	访问方式	xxx:port/项目url/控制器映射键/方法名
//	http://localhost:7080/example/test/QueryList
//	http://localhost:7080/example/test/QueryListMap?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
//	http://localhost:7080/example/test/QueryListWR?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
//	http://localhost:7080/example/test/QueryListStruct?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
//	http://localhost:7080/example/test/QueryListStructWR?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
func TestMVC(t *testing.T){

	//	获得mvc控制器 实例
	mvc := winterMVC.GetInstanceByDispatcherHandler()
	projectPrefix := "example"
	testC := GetInstanceByTestCtrl()
	//	存入控制器
	mvc.RouteCtrl(projectPrefix,"test",&testC)

	//	配置拦截器
	var inter myInterceptor
	mvc.RouteProjectInterceptor(projectPrefix,inter)

	//	配置过滤器
	var f filter
	mvc.SetHttpFilter(f)

	//	错误处理
	var fail failureResp
	mvc.SetFailureResponse(fail)

	//	启动http服务 方式1
	//http.HandleFunc("/winterMvc/", mvc.HandlerFun())
	//http.Handle("/favicon.ico",http.FileServer(http.Dir("./web/img")))
	//http.ListenAndServe(":7080",nil)

	//	启动http服务 方式2
	maxHeaderBytes := 1024 * 1024 * 20
	server := http.Server{
		Addr: ":7080",
		//Handler: mvc,
		ReadTimeout: time.Second * 60,
		WriteTimeout: time.Second * 60,
		IdleTimeout: time.Second * 60,
		MaxHeaderBytes: maxHeaderBytes,
	}
	http.Handle("/favicon.ico",http.FileServer(http.Dir("./web/img")))
	http.Handle("/",mvc)
	fmt.Println(server.ListenAndServe())
}







