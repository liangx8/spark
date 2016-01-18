package spark
import (
	"reflect"
	"net/http"
	"log"

	"github.com/liangx8/spark/invoker"
	"os"
)


type (
	Handler interface{}
	Spark struct{
		invoker.Invoker
		handlers []Handler
		action func(Context,ReturnHandler)
		log log.Logger
	}

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
	c.index ++
	c.run()
}
func (c *context)run(){
	count := len(c.handlers)
	for c.index < count {
		c.Invoke(c.handlers[c.index])
		c.index ++
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
		action:func(c Context,rh ReturnHandler){
			// execute a NotFound response by default
			c.Invoke(rh(http.StatusNotFound,nil))
		},
	}
	spk.action=router.handler
	spk.Invoker=invoker.New()
	spk.Use(Recovery())
	spk.Use(DefaultLogHandler)
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
