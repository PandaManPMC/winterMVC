package winterMVC

import "net/http"

//  author: laoniqiu
//  since: 2022/10/25
//  desc: winterMVC

type ParameterErrorInterface interface {
	ParameterError(http.ResponseWriter, *http.Request, error)
}
