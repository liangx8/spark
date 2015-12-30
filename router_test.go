package spark
import (
	"testing"
	"regexp"
	"net/http"

)

func Test_re(t *testing.T){
	expect := "/xx/bbb/"
	src := expect + "?aaa=bb&cc=(dd)"
	reg := regexp.MustCompile("^(/.*?)\\?")
	ai:=reg.FindStringSubmatchIndex(src)
	if ai == nil {
		t.Fatalf("No match")
	}
	if expect != src[ai[2]:ai[3]] {
		t.Fatalf("expected '%s', but '%s' at %v",expect,src[ai[2]:ai[3]],ai)
	}
}

func Test_router_add(t *testing.T){
	r := &Router{make([]*route,0),[]Handler{http.NotFound}}
	r.Get("/aaa",func(){})
	expectedValue(t,1,len(r.routes))
	expectedValue(t,r.routes[0].match(ANY,"/aaa"),match)
	expectedValue(t,r.routes[0].match(GET,"/aaa?xxxx"),exactMatch)
	expectedValue(t,r.routes[0].match(GET,"/aa"),noMatch)
	expectedValue(t,r.routes[0].match(POST,"/aaa"),noMatch)
	r.Get("/aaa",func(){})
	expectedValue(t,1,len(r.routes))
	r.Any("/aaa",func(){})
	expectedValue(t,2,len(r.routes))
	expectedValue(t,r.routes[1].match(POST,"/aaa"),match)
}

func expectedValue(t *testing.T,i1,i2 interface{}){
	if i1 != i2 {
		t.Fatalf("expected %v of %T, but %v of %T",i1,i1,i2,i2)
	}
}
