package winterMVC

//  author: laoniqiu
//  since: 2022/8/31
//  desc: winterMVC

type LogsInterface interface {
	Info(string)
	Error(string, interface{})
}
