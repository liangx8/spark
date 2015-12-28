package spark

import (
	//	"github.com/liangx8/spark/invoker"
	"net/http"
)
// Distrubutor is a options struct.
// distrubute handler by name string
/*
type Parameters map[string]string
type Distributor struct{
	name string
	data map[string]Handler
}
func NewDistributor(name string)*Distributor{
	return &Distributor{
		name:name,
		data:make(map[string]Handler),
	}
}
func (d *Distributor)Handler(c Context,key string){
	handler,ok := d.data[key]
}
func (d *Distributor)Name()string{
	return d.name
}
*/
type route struct{
	method ReqMethod
	pattern string
	handler []Handler
}
type Router struct{
	routes []*route
	notFound []Handler
}
func (r *Router)Get(p string,h ...Handler)*Router{
	return r.AddRoute(GET,p,h)
}
func (r *Router)Post(p string,h ...Handler)*Router{
	return r.AddRoute(POST,p,h)
}
func (r *Router)Put(p string,h ...Handler)*Router{
	return r.AddRoute(PUT,p,h)
}
func (r *Router)Delete(p string,h ...Handler)*Router{
	return r.AddRoute(DELETE,p,h)
}
func (r *Router)Any(p string,h ...Handler)*Router{
	return r.AddRoute(ANY,p,h)
}
func (r *Router)AddRoute(method ReqMethod,pattern string,h []Handler)*Router{
	r.routes = append(r.routes,&route{method,pattern,h})
	return r
}
func (r *Router)NotFound(h ...Handler)*Router{
	r.notFound = h
	return r
}

func (r *Router)handler(c Context,req *http.Request){

}

type ReqMethod string
const (
	OPTIONS ReqMethod = "OPTIONS"
	GET     ReqMethod = "GET"
	HEAD    ReqMethod = "HEAD"
	POST    ReqMethod = "POST"
	PUT     ReqMethod = "PUT"
	DELETE  ReqMethod = "DELETE"
	TRACE   ReqMethod = "TRACE"
	CONNECT ReqMethod = "CONNECT"
	ANY     ReqMethod = "*"
)
