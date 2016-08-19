package spark

import (
	"net/http"
	"golang.org/x/net/context"
)



type (
	HandleFunc func(context.Context)
	HandleChain func(context.Context,HandleFunc)

	chainLinked struct {
		handle HandleChain
		next *chainLinked
	}
)


func (cl *chainLinked)first(ctx context.Context){
	cl.handle(ctx,cl.moveNext())
}
func (cl *chainLinked)moveNext() HandleFunc{
	var chain HandleFunc
	linked := cl
	chain = func(ctx context.Context){
		if linked.next == nil {
			return // end of chain
		}
		linked = linked.next
		linked.handle(ctx,chain)
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
func wrap(
	cl        *chainLinked,
	handle    HandleFunc,
	f         func(*http.Request)context.Context) http.HandlerFunc{
	cp := cl.Copy()
	chainLinkedAppend(cp,func(ctx context.Context,_ HandleFunc){
		handle(ctx)
	})
	return func(w http.ResponseWriter,r *http.Request){
		cp.first(HttpContext(w,r,f(r)))
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
