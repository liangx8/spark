package spark

import (
	"reflect"
	"net/http"
	"fmt"
)

// statusCode is a type of http status define in package "net/http"
// data is the reflect value of return of action
type ReturnHandler func(statusCode int,data []reflect.Value)Handler

func defaultReturnHandler(statusCode int,data []reflect.Value)Handler{
	if statusCode == http.StatusNotFound {
		return http.NotFound
	}
	count := len(data)
	if count > 0 {
		returnStatus := http.StatusOK
		var returnValue reflect.Value
		if data[0].Kind() == reflect.Int {
			returnStatus = int(data[0].Int())
			if count > 1 {
				returnValue = data[1]
			}
		} else{
			returnValue=data[0]
		}
		if returnStatus == http.StatusNotFound {
			return http.NotFound
		}
		return func(w http.ResponseWriter){
			if returnStatus != http.StatusOK {
				w.WriteHeader(returnStatus)
			}
			fmt.Fprint(w,returnValue)
		}
		

	} 
	// return do nothing if no return
	return func(){}
}
