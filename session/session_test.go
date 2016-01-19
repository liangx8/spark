package session_test
/*
import (
	"testing"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"


	"github.com/liangx8/spark"
	"github.com/liangx8/spark/session"
)

func Test_session(t *testing.T){
	spk := spark.New()
	spk.Use(session.DefaultHandler())
	var sid string
	var s1 session.Session
	spk.GetRouter().Get("/",func(s session.Session) string{
		s1=s
		sid=s.Id()
		if err:=s.Set("name","rose"); err != nil {
			panic(err)
		}
		return s.Id()
	}).Get("/x",func(s session.Session,r *http.Request) string {
		o,err := s.Get("name")
		if err != nil {
			panic(err)
		}
		t.Log("list cookies:")
		for _,cc := range r.Cookies() {
			t.Log(cc)
		}
		expectEqual(t,s1,s)
		expectEqual(t,"rose",o)
		return o.(string)
	})
	ts :=httptest.NewServer(spk)
	defer ts.Close()
	res,err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	cookieStr := res.Header["Set-Cookie"][0]
	cookieId,_ := parse(cookieStr)

	oid := asString(res.Body)
	res.Body.Close()
	expectEqual(t,sid,oid)
	expectEqual(t,sid,cookieId)
//	cookie := &http.Cookie{Name:cookieName,Value:cookieId,Path:"/"}
	
	
	req,err := http.NewRequest("GET",ts.URL+"/x",nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Cookie",cookieStr)

	res, err=http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body,err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}
func parse(str string) (id,name string){
	name = str[:11]
	id = str[12:31]
	return
}
func asString(r io.Reader)string{
	bs,err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(bs)
}
func expectEqual(t *testing.T,o1,o2 interface{}){
	if o1 != o2 {
		t.Fatalf("expected %v of %T, but %v of %T",o1,o1,o2,o2)
	}
}

*/
