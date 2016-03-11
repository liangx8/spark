package spark

import (
	"reflect"
	"net/http"
	"fmt"
)

// statusCode is a type of http status define in package "net/http"
// data is the reflect value of return of action
type (
	ReturnHandlerChain func(int,[]reflect.Value)Handler
	ReturnHandler func(statusCode int,data []reflect.Value,chain ReturnHandlerChain)Handler
	returnHandlerLinked struct {
		rh ReturnHandler
		next *returnHandlerLinked
	}
)
func newReturnHandlerLinked() *returnHandlerLinked{
	return returnHandlerLinkedInsert(nil,defaultReturnHandler)
}
func returnHandlerLinkedInsert(rhl *returnHandlerLinked,rh ReturnHandler) *returnHandlerLinked{
	return &returnHandlerLinked{rh:rh,next:rhl}
}
func (rhl *returnHandlerLinked)First(code int,data []reflect.Value) Handler{
	return rhl.rh(code,data,rhl.Next())
}
func (rhl *returnHandlerLinked)Next() ReturnHandlerChain{
	var chain ReturnHandlerChain
	linked := rhl
	chain = func(statusCode int,data []reflect.Value)Handler{
		if linked.next == nil {
			return doNothing// end of chain.
		}
		next := linked.next
		linked = linked.next // forward the pointer
		return next.rh(statusCode,data,chain)
	}
	return chain
}

func defaultReturnHandler(statusCode int,data []reflect.Value,chain ReturnHandlerChain)Handler{
	if statusCode == http.StatusNotFound {
		// expected a reflect value of string or panic
		return NotFound(data[0].String())
	}
	if statusCode == http.StatusInternalServerError {
		// expecting return value are a single reflect value of error
		return InternalError(data[0].Interface().(error))
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
			res := ""
			if count > 1{
				if data[1].Kind() == reflect.String {
					res = data[1].String()
				}
			}
			return NotFound(res)
		}
		if returnStatus == http.StatusInternalServerError {
			var err error
			if count > 1 {
				var ok bool
				// expected a reflect value of error
				err,ok=data[1].Interface().(error)
				if !ok {
					return InternalError(nil)
				}
			}
			return InternalError(err)
		}
		return func(w http.ResponseWriter){
			if returnStatus != http.StatusOK {
				w.WriteHeader(returnStatus)
			}
			fmt.Fprint(w,returnValue)
		}
		

	} 
	// return do nothing if no return
	return doNothing
}
func doNothing(){}

func InternalError(err error)Handler{
	return func(w http.ResponseWriter,r *http.Request){
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html")
		s:="application error"
		w.Write([]byte(fmt.Sprintf(errorHtml,s,s,err)))
	}
}


const errorHtml = `<html>
<head><title>ERROR: %s</title>
<style type="text/css">
html, body {
font-family: "Roboto", sans-serif;
color: #333333;
background-color: #ea5343;
margin: 0px;
}
h1 {
color: #d04526;
background-color: #ffffff;
padding: 20px;
border-bottom: 1px dashed #2b3848;
}
pre {
margin: 20px;
padding: 20px;
border: 2px solid #2b3848;
background-color: #ffffff;
}
</style>
</head><body>
<h1>ERROR</h1>
<pre style="font-weight: bold;">%s</pre>
<pre>%v</pre>
</body>
</html>`
