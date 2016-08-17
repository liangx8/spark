package session_test

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"fmt"


	"github.com/liangx8/spark/session"
)

func Test_base_session(t *testing.T){

	session.BaseSessionInit()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		s:=session.Get(w,r)
		fmt.Fprint(w,s.Id())
	}))
	defer ts.Close()
	client := &http.Client{}
	req,err := http.NewRequest("GET",ts.URL,nil)
	res,err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body,err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	sessionId:=res.Cookies()[0].Value
	sameOrError(t,body,sessionId)
//	req.AddCookie(res.Cookies()[0])
	res,err = client.Do(req)
	body,err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	sessionId=res.Cookies()[0].Value
	sameOrError(t,body,sessionId)
}
func sameOrError(t *testing.T,b []byte,s string){
	for i,c := range []byte(s) {
		if c != b[i] {
			t.Errorf("Not same")
			return
		}
	}
}
