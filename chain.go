package spark

import (
	"net/http"
)

type (
	HandleFunc func(http.ResponseWriter, *http.Request)
	HandleChain func(w http.ResponseWriter, r *http.Request,chain HandleFunc)

	chainLinked struct {
		handle HandleChain
		next *chainLinked
	}
)

func (cl *chainLinked)first(w http.ResponseWriter,r *http.Request){
	cl.handle(w,r,cl.moveNext())
}
func (cl *chainLinked)moveNext() HandleFunc{
	var chain HandleFunc
	linked := cl
	chain = func(w http.ResponseWriter,r *http.Request){
		if linked.next == nil {
			return // end of chain
		}
		linked = cl.next
		linked.handle(w,r,chain)
	}
	return chain
}
func (cl *chainLinked)Copy() *chainLinked{
	this := cl
	head := &chainLinked{handle:this.handle}
	phead :=head
	for this.next != nil {
		this=this.next
		phead.next=&chainLinked{handle:this.handle}
		phead=phead.next
	}
	return head
}
func (cl *chainLinked)wrap(handle HandleFunc) HandleFunc{
	cp := cl.Copy()
	chainLinkedAppend(cp,func(w http.ResponseWriter,r *http.Request,_ HandleFunc){
		handle(w,r)
	})
	return func(w http.ResponseWriter,r *http.Request){
		cp.first(w,r)
	}
}
func chainLinkedInsert(cl *chainLinked,chain HandleChain) *chainLinked{
	return &chainLinked{handle:chain,next:cl}
}

func chainLinkedAppend(cl *chainLinked, chain HandleChain){
	if cl.next == nil {
		cl.next=&chainLinked{handle:chain,next:nil}
		return
	}
	chainLinkedAppend(cl.next,chain)
}
