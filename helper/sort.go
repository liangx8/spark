package helper

import (
	"sort"
)

type (
	Less func(l,r interface{}) bool
	sortInterface struct{
		ary Array
		less Less
	}
)

func (si *sortInterface) Len() int {
	return si.ary.Len()
}
func (si *sortInterface) Less(i,j int)bool{
	var l,r interface{}
	si.ary.Get(i,&l)
	si.ary.Get(j,&r)
	return si.less(l,r)
}
func (si *sortInterface) Swap(i,j int){
	var l,r interface{}
	si.ary.Get(i,&l)
	si.ary.Get(j,&r)
	si.ary.Set(i,r)
	si.ary.Set(j,l)
	
}
func NewSorter(a interface{}, less Less) sort.Interface{
	return &sortInterface{ary:NewArray(a),less:less}
}
