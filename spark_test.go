package spark_test

import (
	"io"
	"io/ioutil"
	"testing"
	"net/http"
	"fmt"

	"net/http/httptest"
	"github.com/liangx8/spark"
)


func Test_spark(t *testing.T){
	spk := spark.New()
	spk.Use(func(w http.ResponseWriter){
		fmt.Println("user space")
	})
	ts :=httptest.NewServer(spk)
	defer ts.Close()
	res,err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	expected_equal_n(t,res.Body,"404 ",4)
	res.Body.Close()
	
}

func expected_equal(t *testing.T,r io.Reader,expect string){
	greeting,err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if expect != string(greeting) {
		t.Fatalf("expected '%s' but '%s'",expect,greeting)
	}
}
func expected_equal_n(t *testing.T,r io.Reader,expect string,n int){
	buf := make([]byte,n)
	_,err:=r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	if expect[:n] != string(buf) {
		t.Fatalf("expected '%s' but '%s'",expect,buf)
	}
}

