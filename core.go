package spark
import (
	"reflect"
	"net/http"

	"github.com/liangx8/spark/invoker"
	"github.com/liangx8/spark/order"
)
const MAXINT = int((^uint(0)) >> 1)

type (
	Handler interface{}
	Spark struct{
		invoker.Invoker
		handlers order.Order
//		handlers []Handler
//		action func(Context,ReturnHandler)
		GetRouter func()*Router
	}
	// 在Next()方法包裹的中间件中.如果有返回false值.就会放弃剩下的中间件执行,包括action
	Context interface{
		invoker.Invoker
		Next()
	}
	context struct{
		invoker.Invoker
		handlers order.Sequence
//		action func(Context,ReturnHandler)
		chainBreak bool
	}

	responseWriter struct{
		http.ResponseWriter
		Status int
	}
	/*
	nameHandler{
		name string
		handler Handler
	}*/
)
func (c *context)Next(){
	if c.chainBreak { return }
	if c.handlers.Next() {
		c.run()
	}
}
func (c *context)run(){
	for {
		vs:=c.Invoke(c.handlers.Object())
		if len(vs)>0 {

			// break chain if return is false
			if vs[0].Kind() == reflect.Bool && !vs[0].Bool() {
				// break middleware chain
				c.chainBreak=true
				return
			}
		}
		if c.chainBreak { return }
		if !c.handlers.Next() { return }
	}
}
// use a specified seq middleware
func (spk *Spark)UseSeq(seq int,name string,mw Handler){
	check(mw)
	spk.handlers.Add(seq,mw)
}
// mw middleware

// name name of middleware for logging use

func (spk *Spark)Use(name string,mw Handler){
	spk.UseSeq(100,name,mw)
//	spk.handlers = append(spk.handlers,mw)
}
func (spk *Spark)RunAt(addr string){
	http.ListenAndServe(addr,spk)
}
func (spk *Spark)Run(){
	spk.RunAt(":8080")
}
func check(h Handler){
	if reflect.TypeOf(h).Kind() != reflect.Func {
		panic("Handler must be a function")
	}
}
func New() *Spark{
	router:=newRouter()
	
	spk := &Spark{
		handlers:order.New(),
		GetRouter:func()*Router{
			return router
		},
	}
	// add router service to the end of middleware chain

	spk.Invoker=invoker.New()
	spk.Use("logger",DefaultLogHandler)
	spk.Use("recovery",Recovery())
	spk.UseSeq(MAXINT,"action",router.handler)
	spk.Map(DefaultLog())
	spk.Map(ReturnHandler(defaultReturnHandler))
	return spk
}
// implement http.HandlerFunc
func (spk *Spark)ServeHTTP(w http.ResponseWriter,r *http.Request){
	ctx := &context{
		handlers:spk.handlers.NewSequence(),
//		action:spk.action,
		chainBreak:false,
	}
	ctx.Invoker = invoker.New()
	ctx.SetParent(spk)
	ctx.Map(r)
	ctx.MapTo(&responseWriter{w,http.StatusOK},(*http.ResponseWriter)(nil))

	ctx.MapTo(ctx,(*Context)(nil))
	ctx.Next()
}
func (w *responseWriter)WriteHeader(status int){
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}
