package order_test

import (
	"testing"
	odr "github.com/liangx8/spark/order"
)

func Test_order(t *testing.T){
	order := odr.New()
	seq := order.NewSequence()
	if seq.Next() {
		t.Error("expected a empty Order return false")
	}
	order = odr.New()
	order.Add(10,3)
	order.Add(9,1)
	order.Add(10,4)
	order.Add(8,0)
	order.Add(9,2)
	order.Add(10,5)
	
	seq = order.NewSequence()
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
