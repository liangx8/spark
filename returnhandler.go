package spark
import (
	"reflect"
	"net/http"

)

type ReturnHandler func(c Context,retval []reflect.Value)


func DefaultReturnHandler(c Context,retval []reflect.Value){
	wv := c.Get(reflect.TypeOf((*http.ResponseWriter)(nil)).Elem())
	w := wv.Interface().(http.ResponseWriter)
	var responseVal reflect.Value
	if len(retval) > 1 && retval[0].Kind() == reflect.Int {
		w.WriteHeader(int(retval[0].Int()))
		responseVal = retval[1]
	} else if len(retval) > 0 {
		responseVal = retval[0]
	}
	if canDeref(responseVal) {
		responseVal = responseVal.Elem()
	}
	if isByteSlice(responseVal) {
		w.Write(responseVal.Bytes())
	} else {
		w.Write([]byte(responseVal.String()))
	}
	
}

func isByteSlice(val reflect.Value) bool {
	return val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.Uint8
}

func canDeref(val reflect.Value) bool {
	return val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr
}
