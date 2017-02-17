package helper_test

import (
	"testing"
	"github.com/liangx8/spark/helper"
)


func Test_array(t *testing.T){
	ia :=[...]int{0,1,2,3,4}
	ary := helper.NewArray(ia)
	if ary.Len() != 5{
		t.Fatalf("expect 5 but %d",ary.Len())
	}
	for i,v := range ia{
		var intv int
		ary.Get(i,&intv)
		if v != intv {
			t.Fatalf("expecting %d but %d",v,intv)
		}
	}
	sl := make([]int,5)
	ary = helper.NewArray(sl)
	if ary.Len() != 5{
		t.Fatalf("expect 5 but %d",ary.Len())
	}
	for i,_ := range sl{
		ary.Set(i,i)
	}
	for i,v := range sl{
		if i != v {
			t.Fatalf("expecting %d but %d",v,v)
		}
	}
}
