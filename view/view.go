package view

import (
	"io"
	"reflect"
	"bytes"
	"net/http"

	"github.com/liangx8/spark"

)

type (

	View struct {
		makeRender func(RenderMaker) Render
		data interface{}
		contentType string
		statusCode int
	}
)


func Html(name string,data interface{}) *View {
	return &View{
		makeRender:func(maker RenderMaker) Render{
			return maker.ByName(name)
		},
		data:data,
		contentType:"text/html",
		statusCode:http.StatusOK,
	}
}
func ErrorView(err error) *View{
	return &View{
		makeRender:func(maker RenderMaker) Render{
			render := maker.ByName("error")
			if render != nil {
				return render
			}
			return func(io.Writer,interface{}) error{
				// pass current error to chain
				return err
			}
		},
		data:err,
		contentType:"text/html",
		statusCode:http.StatusInternalServerError,
	}
}

func ViewReturnHandler(maker RenderMaker) spark.ReturnHandler{
	return func(statusCode int,data []reflect.Value,chain spark.ReturnHandlerChain) spark.Handler{
		if len(data)>0 {
			v,ok := data[0].Interface().(*View)
			if ok {
				var buf bytes.Buffer
				render:=v.makeRender(maker)
				if err:=render(&buf,v.data); err != nil {
					// a single reflect value of error must be pass to next return value
					return chain(http.StatusInternalServerError,[]reflect.Value{reflect.ValueOf(err)})
				}
				return func(w http.ResponseWriter){
					w.Header().Set("Content-Type",v.contentType)
					if _,err:=io.Copy(w,&buf); err != nil {
						panic(err)
					}
				}
			}
		}
		return chain(statusCode,data)
	}
}
