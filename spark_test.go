package spark_test

import (
	"io"
	"io/ioutil"
	"testing"
	"net/http"


	"net/http/httptest"
	"github.com/liangx8/spark"
)


func Test_use(t *testing.T){
	spk := spark.New()
	spk.Use("",func(w http.ResponseWriter){
		
	})
	ts :=httptest.NewServer(spk)

	res,err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatal("expected status code 404")
	}
	res.Body.Close()
	ts.Close()
	spk = spark.New()
	spk.Use("",func(w http.ResponseWriter){
		panic("painc")
	})
	ts =httptest.NewServer(spk)

	res,err = http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusInternalServerError {
		t.Fatal("expected status code 500")
	}
	res.Body.Close()
	ts.Close()
	
}

func Test_router(t *testing.T){
	spk := spark.New()
	spk.GetRouter().Get("/x",func()string{
		return "ok"
	})
	spk.GetRouter().Get("/y",func()int {
		return http.StatusNotFound
	})
	spk.GetRouter().Get("/z",func() {
		panic("panic")
	})
	
	ts :=httptest.NewServer(spk)

	res,err := http.Get(ts.URL+"/x")
	if err != nil {
		t.Fatal(err)
	}
	expected_equal(t,res.Body,"ok")
	res,err = http.Get(ts.URL + "/y")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatal("expected status code 404")
	}
	res,err = http.Get(ts.URL + "/z")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusInternalServerError {
		t.Fatal("expected status code 500")
	}
		
	ts.Close()
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

