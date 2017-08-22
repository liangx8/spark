package spark_test

import (

	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"fmt"

	"golang.org/x/net/context"
	"github.com/liangx8/spark"

)
func newTestSpark() *spark.Spark{
	return spark.New(func(_ *http.Request)context.Context{
		return context.Background()
	})
}
func Test_spark(t *testing.T){
	
	spk := newTestSpark()
	spk.AddChain(func(ctx context.Context,chain spark.HandleFunc){
		w,_,err := spark.ReadHttpContext(ctx)
		if err != nil {
			t.Error(err)
		}
		t.Log("chain 2\n")

		fmt.Fprint(w,data[6:12])
		chain(ctx)
	})
	spk.AddChain(func(ctx context.Context,chain spark.HandleFunc){
		w,_,err := spark.ReadHttpContext(ctx)
		if err != nil {
			t.Error(err)
		}
		t.Log("chain 1\n")
		fmt.Fprint(w,data[:6])
		chain(ctx)
	})
	ts := httptest.NewServer(spk.Handler(func(ctx context.Context){
		w,_,err := spark.ReadHttpContext(ctx)
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(w,data[12:])
	}))
	defer ts.Close()
	client := &http.Client{}
	req,err := http.NewRequest("GET",ts.URL,nil)
	if err != nil {	t.Fatal(err)}
	res,err := client.Do(req)
	if err != nil {	t.Fatal(err)}
	body,err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {	t.Fatal(err)}
	sameOrError(t,body,data)
}
func sameOrError(t *testing.T,b []byte,d string){
	for i,c := range []byte(d) {
		if c != b[i] {
			t.Fatal("not same")
		}
	}
}
const data = `this is a test data, 
adkj;lkfj'llkj;lkdjf;lkasjd;lkfj akelfdjkqwer
adfklq qwerjkj;lq qwerqwer qewr  rqwerqker  weqrqwer
jkdjf adsffasdf 

`
