package helper_test

import (
	"testing"
	"github.com/liangx8/spark/helper"
	"sort"
)


func Test_sort(t *testing.T){
	sl := make([]int,5)
	sl[0]=4
	sl[1]=2
	sl[2]=3
	sl[3]=1
	sl[4]=0
	sor := helper.NewSorter(sl,func(l,r interface{}) bool{
		return l.(int)<r.(int)
	})
	sort.Sort(sor)

	for i,v := range sl{
		if i != v {
			t.Fatalf("expecting %d(%T) but %d(%T)",i,i,v,v)
		}
	}

}
