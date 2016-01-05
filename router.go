package spark

import (
	"net/http"
	"regexp"
	"reflect"

	"github.com/liangx8/spark/invoker"
)

type(
	route struct{
		method ReqMethod
		url string
		handlers []Handler
	}
	Router struct{
		routes []*route
		notFound []Handler
	}
	routeContext struct{
		invoker.Invoker
		handlers []Handler
		index int
		returnHandler ReturnHandler
	}
	routeMatch int
	ReqMethod string
)
var urlReg = regexp.MustCompile("^(.*?)\\?")
func (r *route)match(method ReqMethod,url string) routeMatch {
	idx:=urlReg.FindStringSubmatchIndex(url)
	var prefix string
	if idx == nil {
		prefix=url
	} else {
		prefix=url[idx[2]:idx[3]]
	}
	if prefix == r.url {
		if method == r.method {
			return exactMatch
		}
		if method == ANY {
			return match
		}
		if ANY == r.method {
			return match
		}
	}
	return noMatch
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
func (r *Router)AddRoute(method ReqMethod,url string,h []Handler)*Router{
	for i,ro := range r.routes {
		if m :=ro.match(method,url);m == exactMatch {
			r.routes[i]=&route{method,url,h}
			return r
		}
	}
	r.routes = append(r.routes,&route{method,url,h})
	return r
}
func (r *Router)NotFound(h ...Handler)*Router{
	r.notFound = h
	return r
}

func (r *Router)handler(ctx Context,req *http.Request,rh ReturnHandler){

	var rt *route
	var hs []Handler
	mm := noMatch
	for _,ro := range r.routes {
		m := ro.match(ReqMethod(req.Method),req.URL.Path)
		if m == exactMatch{
			rt = ro
			mm = m
			break
		}
		if m > mm {
			rt = ro
			mm = m
		}
	}
	if mm == noMatch {
		hs = r.notFound[:]
	} else {
		hs = rt.handlers[:]
	}
	c := &routeContext{invoker.New(),hs,0,nil}
	c.SetParent(ctx)
	if rh == (ReturnHandler)(nil) {
		c.returnHandler=DefaultReturnHandler
	} else {
		c.returnHandler=rh
	}
	c.MapTo(c,(*Context)(nil))
	c.run()
}
func (c *routeContext)OnReturn(vals []reflect.Value){
	if c.returnHandler != nil {
		c.returnHandler(c,vals)
	}
}
func (c *routeContext)Next(){
	c.index ++
	c.run()
}
func (c *routeContext)run(){
	for c.index < len(c.handlers) {
		vals := c.Invoke(c.handlers[c.index])
		c.OnReturn(vals)
		c.index ++
	}
}

const (
	noMatch routeMatch = iota
	match
	exactMatch
)
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
