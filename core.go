package spark

import (
	"net/http"
)

type (
	Spark struct{
		head *chainLinked
	}
)
// HandleChain is invoked in revse order
func (spk *Spark)AddChain(chain HandleChain){
	if spk.head == nil {
		spk.head = chainLinkedInsert(nil,chain)
	} else {
		spk.head = chainLinkedInsert(spk.head,chain)
	}
}

func (spk *Spark)HandleFunc(pattern string, handle HandleFunc){
	if spk.head != nil {
		http.HandleFunc(pattern,spk.head.wrap(handle))
	} else {
		http.HandleFunc(pattern,handle)
	}
}
