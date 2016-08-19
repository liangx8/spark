package session_test

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"fmt"

	"golang.org/x/net/context"

	"github.com/liangx8/spark"
	"github.com/liangx8/spark/session"
)

func Test_base_session(t *testing.T){

	spk :=spark.New(func(_ *http.Request)context.Context{
		return context.Background()
	})
	spk.AddChain(session.CreateSessionChain())
	ts := httptest.NewServer(spk.Handler(func(ctx context.Context){
		s,err:=session.GetSession(ctx)
		if err != nil {
			t.Error(err)
		}
		w,_,err := spark.ReadHttpContext(ctx)
		if err != nil {
			t.Error(err)
		}
		t.Log(s.Id())
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
	req.AddCookie(res.Cookies()[0])
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
const data = `this is a test data, 
adkj;lkfj'llkj;lkdjf;lkasjd;lkfj akelfdjkqwer
adfklq qwerjkj;lq qwerqwer qewr  rqwerqker  weqrqwer
jkdjf adsffasdf 

`
