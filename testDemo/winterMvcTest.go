package testDemo

import (
	"../winter-mvc-core"
	"net/http"
	"time"
)


//	测试
//	访问方式	xxx:port/项目url/控制器映射键/方法名
//	http://localhost:18080/winterMvc/test/QueryList
//	http://localhost:18080/winterMvc/test/QueryList2?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
//	http://localhost:18080/winterMvc/test/QueryList3?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
//	http://localhost:18080/winterMvc/test/QueryList4?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
//	http://localhost:18080/winterMvc/test/QueryList5?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
func TestWinterMvc(){

	//	获得mvc控制器 实例
	mvc := winter_mvc_core.GetInstanceByDispatcherHandler()
	projectPrefix := "winterMvc"
	ctrlPrefix := "test"

	testC := GetInstanceByTestCtrl()
	//	存入控制器
	mvc.PutCtrlPrefix(projectPrefix,ctrlPrefix,&testC)

	//	配置拦截器
	var inter myInterceptor
	mvc.PutProjectInterceptor(projectPrefix,inter)

	//	配置过滤器
	var f filter
	mvc.SetHttpFilter(f)

	//	错误处理
	var fail failureResp
	mvc.SetFailureResponse(fail)

	//	启动http服务 方式1
	//http.HandleFunc("/winterMvc/", mvc.HandlerFun())
	//http.Handle("/favicon.ico",http.FileServer(http.Dir("./web/img")))
	//http.ListenAndServe(":18080",nil)


	//	启动http服务 方式2
	maxHeaderBytes := 1024 * 1024 * 20
	server := http.Server{
		Addr: ":18080",
		//Handler: mvc,
		ReadTimeout: time.Second * 60,
		WriteTimeout: time.Second * 60,
		IdleTimeout: time.Second * 60,
		MaxHeaderBytes: maxHeaderBytes,
	}
	http.Handle("/favicon.ico",http.FileServer(http.Dir("./web/img")))
	http.Handle("/",mvc)
	server.ListenAndServe()
}







