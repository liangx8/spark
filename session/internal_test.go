package session

import (
	"testing"
	"log"

	"os"
)

func Test_session(t *testing.T){
	l:=log.New(os.Stdout,"testing ",log.Ltime)
	sc := make(sessionContext)
	s := sc.getOrNew("",l)
	oldId := s.Id()
	expectEqual(t,s,sc.getOrNew(oldId,l))
	s.Set("one",1)
	expectEqual(t,1,sc.getOrNew(oldId,l).Get("one"))
	s.Invalidate()
	expectUnequal(t,s,sc.getOrNew(oldId,l))
	
}

func expectUnequal(t *testing.T, o1,o2 interface{}){
	if o1 == o2 {
		t.Fatalf("expected %v of %T is not equal",o1)
	}
}
func expectEqual(t *testing.T,o1,o2 interface{}){
	if o1 != o2 {
		t.Fatalf("expected '%v' of %T, but '%v' of %T",o1,o1,o2,o2)
	}
}

