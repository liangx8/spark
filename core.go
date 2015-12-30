package spark
import (
	"reflect"
	"net/http"


	"github.com/liangx8/spark/invoker"
)
// Both middleware and action use this type

type Handler interface{}

func validate(h Handler){
	if reflect.TypeOf(h).Kind() != reflect.Func {
		panic("Handler must be a function")
	}
}
type Spark struct{
	invoker.Invoker
	handlers []Handler
	action Handler
	log Logger
	GetRouter func()*Router
}

func New() *Spark{
	r := &Router{make([]*route,0),[]Handler{http.NotFound}}
	
	spk := &Spark{
		invoker.New(),
		make([]Handler,0),
		r.handler,
		DefaultLogger(),
		func()*Router{return r},
	}
	spk.MapTo(spk.log,(*Logger)(nil))
	spk.Use(Recovery())
	return spk
}
func (spk *Spark)Use(h Handler){
	validate(h)
	spk.handlers = append(spk.handlers,h)
}
func (spk *Spark)ServeHTTP(w http.ResponseWriter,r *http.Request){

	c:=&context{invoker.New(),spk.handlers,spk.action,0,nil}
	c.Map(w)
	c.Map(r)
	c.SetParent(spk)
	c.MapTo(c,(*Context)(nil))
	c.run()
}


type Context interface{
	invoker.Invoker
	Next()
}
type context struct{
	invoker.Invoker
	handlers []Handler
	action Handler
	index int
	returnHandler ReturnHandler
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
			if c.returnHandler != nil {
				c.returnHandler(c,vals)
			}
		}
		return
	}
	if c.index > cnt {
		panic("Never reach here")
	}
}
var DoNothing Handler = func(){}
