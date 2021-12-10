package winter_mvc_core

import "net/http"

type FailureResponseInterface interface {

	Failure404(writer http.ResponseWriter,request *http.Request)

	Failure500(writer http.ResponseWriter,request *http.Request)

}


