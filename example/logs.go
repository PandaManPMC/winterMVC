package example

import "fmt"

//  author: laoniqiu
//  since: 2022/8/31
//  desc: example

type logs struct {
}

func (*logs) Info(msg string) {
	fmt.Println(msg)
}

func (*logs) Error(msg string, err interface{}) {
	fmt.Println(msg)
	fmt.Println(err)
}
