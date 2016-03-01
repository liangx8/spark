package spark

import (
	"testing"
	"reflect"
)

func Test_returnhandler(t *testing.T){
	linked := newReturnHandlerLinked()
	aout := make([]int,2)
	ain := []int{2,0}
	
	linked=returnHandlerLinkedInsert(linked,func(code int,data []reflect.Value,chain ReturnHandlerChain)Handler{
		aout[0]=code
		return chain(1,nil)
	})
	linked=returnHandlerLinkedInsert(linked,func(code int,data []reflect.Value,chain ReturnHandlerChain)Handler{
		aout[1]=code
		return chain(2,nil)
	})
	linked.First(0,nil)
	if ain[0] != aout[0] {
		t.Errorf("expected pos 0 is '%d' but '%d'",ain[0],aout[0])
	}
	if ain[1] != aout[1] {
		t.Errorf("expected pos 0 is '%d' but '%d'",ain[1],aout[1])
	}
}
