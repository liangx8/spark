package spark
import (
	"reflect"
	"net/http"

	"github.com/liangx8/spark/invoker"
	"github.com/liangx8/spark/order"

	netctx "golang.org/x/net/context"
)
const MAXINT = int((^uint(0)) >> 1)

type (
	Handler interface{}
	Spark struct{
		invoker.Invoker
		background netctx.Context
		handlers order.Order
//		handlers []Handler
//		action func(Context,ReturnHandler)
		GetRouter func()*Router
		rhLinked *returnHandlerLinked
		log LogAdaptor
	}
	Context interface{
		invoker.Invoker
		Next()
		Die()
	}
	context struct{
		netctx.Context
		invoker.Invoker
		handlers order.Sequence
		cancel func()
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
func (c *context)Die(){
	c.cancel()
}
func (c *context)Next(){
	c.run()
}

func (c *context)run(){


	for {

		if !c.handlers.Next() {
			c.Die()
			return
		}
		select {
		case <-c.Done():
			return
		default:
			c.Invoke(c.handlers.Object())
		}

	}
}
// use a specified seq middleware
func (spk *Spark)UseSeq(seq int,name string,mw Handler){
	check(mw)
	spk.log(INFO,"Use middleware %s at seqence %d",name,seq)
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
		background:netctx.Background(),
		handlers:order.New(),
		GetRouter:func()*Router{
			return router
		},
		log:DefaultLog(),
		rhLinked:newReturnHandlerLinked(),
	}
	// add router service to the end of middleware chain

	spk.Invoker=invoker.New()
	spk.Use("logger",DefaultLogHandler)
	spk.Use("recovery",Recovery())
	spk.UseSeq(MAXINT,"action",router.handler)
	spk.Map(spk.log)
	spk.Map(spk.rhLinked)
	return spk
}
func (spk *Spark)Start(){
	spk.StartAt(":8080")
}
func (spk *Spark)StartAt(port string){
	spk.log(INFO,"Server starting at port %s",port)
	spk.log(ERROR,"starting error",http.ListenAndServe(port,spk))
}
func (spk *Spark)RegisterReturnHandler(rh ReturnHandler){
	spk.rhLinked=returnHandlerLinkedInsert(spk.rhLinked,rh)
	spk.Map(spk.rhLinked)
}
// implement http.HandlerFunc
func (spk *Spark)ServeHTTP(w http.ResponseWriter,r *http.Request){
	nctx,ctxCancle := netctx.WithCancel(spk.background)
	ctx := &context{
		handlers:spk.handlers.NewSequence(),
		cancel:ctxCancle,
	}
	ctx.Context=nctx
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
