package spark
import (
	"reflect"
	"net/http"
	"log"

	"github.com/liangx8/spark/invoker"
	"os"
)
const MAXINT = int((^uint(0)) >> 1)

type (
	Handler interface{}
	Spark struct{
		invoker.Invoker
		handlers []Handler
		action func(Context,ReturnHandler)
		GetRouter func()*Router
	}
	// 在Next()方法包裹的中间件中.如果有返回false值.就会放弃剩下的中间件执行,包括action
	Context interface{
		invoker.Invoker
		Next()
	}
	context struct{
		invoker.Invoker
		handlers []Handler
		action func(Context,ReturnHandler)
		index int
	}

	responseWriter struct{
		http.ResponseWriter
		Status int
	}
)
func (c *context)Next(){
	if c.index < MAXINT{
		c.index ++
	}
	c.run()
}
func (c *context)run(){
	count := len(c.handlers)
	for c.index < count {
		vs:=c.Invoke(c.handlers[c.index])
		if len(vs)>0 {

			// break chain if return is false
			if vs[0].Kind() == reflect.Bool && !vs[0].Bool() {
				// set a enough large number to prevent chain continue
				c.index= MAXINT
			}
		}
		if c.index < MAXINT{
			c.index ++
		}

	}

	if c.index > count {return}
	c.Invoke(c.action)
}
// use a middleware
func (spk *Spark)Use(mw Handler){
	check(mw)
	spk.handlers = append(spk.handlers,mw)
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
		handlers:make([]Handler,0),
		GetRouter:func()*Router{
			return router
		},
	}
	spk.action=router.handler
	spk.Invoker=invoker.New()
	spk.Use(DefaultLogHandler)
	spk.Use(Recovery())
	spk.Map(log.New(os.Stdout,"[spark] ",log.LstdFlags))
	spk.Map(ReturnHandler(defaultReturnHandler))
	return spk
}
// implement http.HandlerFunc
func (spk *Spark)ServeHTTP(w http.ResponseWriter,r *http.Request){
	ctx := &context{
		handlers:spk.handlers,
		action:spk.action,
		index:0,
	}
	ctx.Invoker = invoker.New()
	ctx.SetParent(spk)
	ctx.Map(r)
	ctx.MapTo(&responseWriter{w,http.StatusOK},(*http.ResponseWriter)(nil))

	ctx.MapTo(ctx,(*Context)(nil))
	ctx.run()
}
func (w *responseWriter)WriteHeader(status int){
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}
