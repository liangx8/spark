package view

import (
	"html/template"
	"reflect"
	"bytes"
	"net/http"

	"github.com/liangx8/spark"
)

type (
	// 还要有 Set-Content 的设置，再想一下怎么设计
	View struct {
		makeRender func(RenderMaker) Render
		data interface{}
		contentType string
	}
)


func Html(name string,data interface{}) *View {
	return &View{
		makeRender:func(maker RenderMaker) Render{
			return maker.ByName(name)
		},
		data:data,
		contentType:"text/html",
	}
}

func ViewReturnHandler(maker RenderMaker) ReturnHandler{
	return func(statusCode int,data []reflect.Value,chain spark.ReturnHandlerChain) spark.Handler{
		if len(data)>0 {
			v,ok := data[0].Interface().(*View)
			if ok {
				var buf bytes.Buffer
				render:=v.makeRender(maker)
				if err:=render(&buf,v.data); err != nil {
					chain(http.StatusInternalServerError,[]reflect.Value{reflect.ValueOf(err)})
				}
				return func(w http.ResponseWriter){
					w.Header().Set("Content-Type",v.contentType)
					if err:=io.Copy(w,&buf); err != nil {
						panic(err)
					}
				}
			}
		}
		return chain(statusCode,data)
	}
}
