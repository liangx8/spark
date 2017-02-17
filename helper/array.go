package helper

import (
	"reflect"
)

type (
	Array interface {
		Len() int
		Get(int, interface{})
		Set(int, interface{})
	}
	arrayImp struct{
		data reflect.Value
	}
)

func (ary *arrayImp)Len() int{
	return ary.data.Len()
}
func (ary *arrayImp)Set(idx int,e interface{}){
	v:=ary.data.Index(idx)
	v.Set(reflect.ValueOf(e))
}
func (ary *arrayImp)Get(idx int,pe interface{}){
	pv := reflect.ValueOf(pe)
	pv.Elem().Set(ary.data.Index(idx))
}

// New Array
func NewArray(ary interface{}) Array{
	return &arrayImp{data:reflect.ValueOf(ary)}
}
