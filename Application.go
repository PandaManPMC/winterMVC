package main

import "./testDemo"

func main() {
	//	测试
	//	访问方式	xxx:port/项目url/控制器映射键/方法名
	//	http://localhost:18080/winterMvc/test/QueryList
	//	http://localhost:18080/winterMvc/test/QueryList2?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
	//	http://localhost:18080/winterMvc/test/QueryList3?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
	//	http://localhost:18080/winterMvc/test/QueryList4?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
	//	http://localhost:18080/winterMvc/test/QueryList5?name=heih&age=99&activity=true&fighting=33.55&inDate=2021-12-09%2014:10:55
	testDemo.TestWinterMvc()
}

