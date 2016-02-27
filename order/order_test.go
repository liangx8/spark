package order_test

import (
	"testing"
	"github.com/liangx8/spark/order"
)

func Test_order(t *testing.T){
	odr := order.New()
	seq := odr.NewSequence()
	if seq.Next() {
		t.Error("expected a empty Order return false")
	}

	odr.Add(10,3)
	odr.Add(9,1)
	odr.Add(10,4)
	odr.Add(8,0)
	odr.Add(9,2)
	odr.Add(10,5)
	
	seq = odr.NewSequence()
	target := 0
	if !seq.Next() {
		t.Error("expected true at first run")
	}
	for {
		obj := seq.Object()
		if target != obj.(int) {
			t.Errorf("expected %d,%v",target,obj)
		}
		target ++
		if !seq.Next() {break}
	}
	
}
