package view
import (
	"fmt"
	"sync"
	"net/http"
	"github.com/liangx8/spark"
	"reflect"
)
type (
	View struct{
		Render spark.Handler
	}
)
var (
	once sync.Once
	maker RenderMaker
)

func ViewReturnHandler(cfg *Config) spark.ReturnHandler{
	m,err := DefaultRenderMaker(cfg)
	if err != nil {
		panic(err)
	}
	return ViewReturnHandlerRenderMaker(m)
}
func ViewReturnHandlerRenderMaker(m RenderMaker) spark.ReturnHandler{
	once.Do(func(){
		if m == nil {
			panic("RenderMaker is not allow")
		}
		maker=m
	})
	return func(c spark.Context,vals []reflect.Value){
		var v *View
		var reterr error
		for _,val := range vals {
			if val.Kind() == reflect.Interface {
				if val.Type().Name() == "error" {
					if !val.IsNil(){
						reterr = val.Interface().(error)
					}
				}
			} else {
				v = val.Interface().(*View)
			}
		}
		if reterr != nil {
			c.Invoke(internalError(reterr))
			return
		}
		if v != nil {
			c.Invoke(v.Render)
		}
	}
}
func Html(name string,data interface{}) *View{
	r :=maker.ByName(name)
	if r == nil {
		panic(fmt.Errorf("template %s is not found",name))
	}
	return &View{
		Render:func(w http.ResponseWriter,log spark.Logger){
			log.Infof("Render correct")
			before := func()error{
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusInternalServerError)
				return nil
			}
			if err:=r(w,before,data);err != nil {

				maker.ByName(errorStr)(w,before,map[string]interface{}{
					"type":"Internal Server Error",
					"error":err,
				})
			}
		},
	}
}
func internalError(err error)spark.Handler{
	return func(w http.ResponseWriter){
		before := func()error{
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)
			return nil
		}
		maker.ByName(errorStr)(w,before,map[string]interface{}{
			"type":"Internal Server Error",
			"error":err,
		})
	}
}
func NotFound(w http.ResponseWriter,r *http.Request){
	w.WriteHeader(http.StatusNotFound)
	maker.ByName(notfoundStr)(
		w,func()error{
			w.Header().Set("Content-Type", "text/html")
			return nil
		},
		map[string]interface{}{
			"type":"Resource not found",
			"error":fmt.Sprintf("%s is not available on this site",r.URL.Path),
		},
	)
}

