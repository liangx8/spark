package spark

import (
	"net/http"
	"errors"
	"golang.org/x/net/context"
)

type (
	Spark struct{
		head *chainLinked
		newContext func(*http.Request)context.Context
	}
	httpPair struct{
		w http.ResponseWriter
		r *http.Request
	}
	typ int
)
const httpPairKey typ = 1

func HttpContext(w http.ResponseWriter,r *http.Request,parent context.Context)context.Context{
	return context.WithValue(parent,httpPairKey,&httpPair{w:w,r:r})
}
func ReadHttpContext(c context.Context) (http.ResponseWriter,*http.Request,error){
	a,ok := c.Value(httpPairKey).(*httpPair)
	if !ok {
		return nil,nil,ErrNotHttpContext
	}
	return a.w,a.r,nil
}
func New(n func(*http.Request)context.Context) *Spark{
	return &Spark{newContext:n}
}
// HandleChain is invoked in revsed order
func (spk *Spark)AddChain(chain HandleChain){
	if spk.head == nil {
		spk.head = chainLinkedInsert(nil,chain)
	} else {
		spk.head = chainLinkedInsert(spk.head,chain)
	}
}

func (spk *Spark)HandleFunc(pattern string, handle HandleFunc){
	if spk.head != nil {
		http.HandleFunc(pattern,wrap(spk.head,handle,spk.newContext))
	} else {
		http.HandleFunc(pattern,convert(handle,spk.newContext))
	}
}

func (spk *Spark)Handler(f HandleFunc)http.HandlerFunc{
	if spk.head != nil {
		return wrap(spk.head,f,spk.newContext)

	} else {
		return convert(f,spk.newContext)

	}
}
func convert(f HandleFunc,nc func(*http.Request)context.Context)http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		f(HttpContext(w,r,nc(r)))
	}
}

var (
	ErrNotHttpContext=errors.New("Not http context")
)
