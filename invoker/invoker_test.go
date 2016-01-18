package invoker_test

import (
	"testing"

	
	"github.com/liangx8/spark/invoker"
)
type iface interface{
	Int() int
}
type struc struct{}
func (struc)Int() int {	return 100 }
type struc2 struct{
	struc
}
func (struc2)Int() int { return 200 }
	

func Test_invoker(t *testing.T){
	ink := invoker.New()
	
	ink.Map(&struc{})
	ink.Invoke(func(s1 *struc,s2 iface){
		expectedValue(t,s1,s2)
	})
	ink = invoker.New()
	ps :=&struc{}
	ink.MapTo(ps,(*iface)(nil))
	ink.Invoke(func(s1 *struc,s2 iface){
		expectedValue(t,ps,s2)
		expectedValue(t,(*struc)(nil),s1)
	})
}
type Num int
func Test_builintype(t *testing.T){
	ink := invoker.New()
	ink.Map(20)
	var ff float32= 210.0
	ink.Map(ff)
	ink.Invoke(func(i int,f float32,name string){
		expectedValue(t,20,i)
		expectedValue(t,ff,f)
		expectedValue(t,"",name)
	})
	n := "name"
	ink.Map(n)
	ink.Invoke(func(name string){
		expectedValue(t,n,name)
	})
	f := func(n Num,ii int){
		expectedValue(t,Num(0),n)
		expectedValue(t,20,ii)
	}
	ink.Invoke(f)
	ink1:=invoker.New()
	ink1.SetParent(ink)
	ink1.Invoke(f)
	ink = invoker.New()
}

func expectedValue(t *testing.T,expected,but interface{}){
	if expected != but {
		t.Fatalf("Expected %v with type %T, but %v with type %T",expected,expected,but,but)
	}
}
