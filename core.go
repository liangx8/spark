package spark
import (
	"reflect"
	"net/http"


	"github.com/liangx8/spark/invoker"
)
// Both middleware and action use this type

type (
	Handler interface{}

	Spark struct{
		invoker.Invoker
		handlers []Handler
		action Handler
		log Logger
		GetRouter func()*Router
	}

	Context interface{
		invoker.Invoker
		Next()
		OnReturn([] reflect.Value)
	}
	context struct{
		invoker.Invoker
		handlers []Handler
		action Handler
		index int
		returnHandler ReturnHandler
	}

	ResponseWriter struct{
		http.ResponseWriter
		Status int
	}
)
func validate(h Handler){
	if reflect.TypeOf(h).Kind() != reflect.Func {
		panic("Handler must be a function")
	}
}

func New() *Spark{
	r := &Router{make([]*route,0),[]Handler{NotFound}}
	
	spk := &Spark{
		handlers:make([]Handler,0),
		action:r.handler,
		log:DefaultLogger(),
		GetRouter:func()*Router{return r},
	}
	spk.Invoker=invoker.New()
	spk.MapTo(spk.log,(*Logger)(nil))
	spk.Use(DefaultLogHandler)
	spk.Use(Recovery())
	return spk
}
func (spk *Spark)Use(h Handler){
	validate(h)
	spk.handlers = append(spk.handlers,h)
}
func (spk *Spark)ServeHTTP(w http.ResponseWriter,r *http.Request){

	c:=&context{invoker.New(),spk.handlers,spk.action,0,nil}
	c.MapTo(&ResponseWriter{w,http.StatusOK},(*http.ResponseWriter)(nil))
	c.Map(r)
	c.SetParent(spk)
	c.MapTo(c,(*Context)(nil))
	c.run()
}
func (spk *Spark)RunAt(addr string){
	http.ListenAndServe(addr,spk)
}
func (spk *Spark)Run(){
	spk.RunAt(":8080")
}
func (c *context)OnReturn(vals []reflect.Value){
	if c.returnHandler != nil {
		c.returnHandler(c,vals)
	}
}
func (c *context)Next(){
	c.index ++
	c.run()
}
func (c *context)run(){
	cnt := len(c.handlers)

	for c.index < cnt {
		vals:=c.Invoke(c.handlers[c.index])
		if c.returnHandler != nil {
			c.returnHandler(c,vals)
		}
		c.index ++
		if c.index > cnt {
			return
		}
	}
	if c.index == cnt {

		vals:=c.Invoke(c.action)
		if len(vals) > 0 {
			c.OnReturn(vals)
		}
		return
	}
	if c.index > cnt {
		panic("Never reach here")
	}
}
func (w *ResponseWriter)WriteHeader(status int){
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}
var DoNothing Handler = func(){}
