# winterMvc
    基于go http 封装 简洁的web服务器框架
    无须任何外部依赖，只需导入 winterMVC    

    go get github.com/PandaManPMC/winterMVC@v1.0.2

### 测试
    
    运行 Application
    request_test.go 运行测试
    或者
    访问连接
    
    //	访问方式	xxx:port/项目url/控制器映射键/方法名
    //	http://localhost:7080/example/test/QueryList
    //	http://localhost:7080/example/test/QueryListMap?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
    //	http://localhost:7080/example/test/QueryListWR?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
    //	http://localhost:7080/example/test/QueryListStruct?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
    //	http://localhost:7080/example/test/QueryListStructWR?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55

### 配置

    import "github.com/PandaManPMC/winterMVC"
    mvc := winterMVC.GetInstanceByDispatcherHandler()

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

### 启动服务

	启动http服务 方式1
	http.HandleFunc("/winterMvc/", mvc.HandlerFun())
	http.Handle("/favicon.ico",http.FileServer(http.Dir("./web/img")))
	http.ListenAndServe(":7080",nil)


	启动http服务 方式2
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

### 控制器编写规则，参照testDemo下testCtrl

    type testCtrl struct {
    }

    //  不接收任何参数
    func (t *testCtrl) QueryList() result{}
    //  接收map[string]string参数
    func (t *testCtrl) QueryList2(params map[string]string) result{}
    //  接收 w http.ResponseWriter, request *http.Request
    func (t *testCtrl) QueryList3(w http.ResponseWriter, request *http.Request) result{}
    //  接收 一个struct
    func (t *testCtrl) QueryList4(t structType) result{}
    //  接收三个参数 struct， w http.ResponseWriter, request *http.Request
    func (t *testCtrl) QueryList5(t structType,w http.ResponseWriter, request *http.Request) result{}


### 过滤器 HttpFilterInterface

    //	过滤
    //	writer http.ResponseWriter
    //	request *http.Request
    //	bool	true继续执行，false中断本次请求
    Filter(writer http.ResponseWriter, request *http.Request) bool


### 拦截器 HandlerInterceptorInterface

    //	在处理器处理请求之前
    //	writer http.ResponseWriter
    //	request *http.Request
    //	bool	响应 true 执行处理，false不继续执行处理也不会执行AfterHandler()
    //	string	在bool为false拦截下请求时的响应数据
    BeforeHandler(writer http.ResponseWriter, request *http.Request) (bool,string)

    //	处理器处理请求之后
    //	writer http.ResponseWriter
    //	request *http.Request
    AfterHandler(writer http.ResponseWriter, request *http.Request)


### 错误处理 FailureResponseInterface

	//	错误处理回调404
	//	找不到对应控制器或处理方法回调404
	Failure404(writer http.ResponseWriter,request *http.Request)

	//	错误处理回调500
	//	出现预料之外的事，如url截取后长度不足3位
	Failure500(writer http.ResponseWriter,request *http.Request)










