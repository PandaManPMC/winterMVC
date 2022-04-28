package example

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T){
	resp,err := http.Get("http://localhost:7080/example/test/QueryList")
	fmt.Println(err)
	fmt.Println("header:",resp.Header)
	body := resp.Body
	buf := make([]byte,1024)
	var str string
	for {
		n,re := body.Read(buf)
		if io.EOF == re && 0 == n{
			fmt.Println("读取完成")
			body.Close()
			break
		}

		if 0 < n{
			str = fmt.Sprintf("%s%s",str,string(buf[0:n]))
		}
	}
	fmt.Println(str)
}





