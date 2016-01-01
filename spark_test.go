package spark_test

import (
	"io"
	"io/ioutil"
	"testing"
	"net/http"
//	"fmt"

	"net/http/httptest"
	"github.com/liangx8/spark"
)


func Test_spark(t *testing.T){
	spk := spark.New()
	spk.Use(func(w http.ResponseWriter){
//		fmt.Println("user space")
	})
	spk.GetRouter().Get("/",func()string{
		return "ok"
	}).Get("/x",func()string{
		return "x-ok"
	}).Get("/y",func() map[string]string{
		return map[string]string{
			"name":"rose",
		}
	})
	ts :=httptest.NewServer(spk)
	defer ts.Close()
	res,err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	expected_ok(t,res.Body,"ok")
	res.Body.Close()

	res,err = http.Get(ts.URL+"/x")
	if err != nil {
		t.Fatal(err)
	}
	expected_ok(t,res.Body,"x-ok")
	res.Body.Close()

	res,err = http.Get(ts.URL+"/y")
	if err != nil {
		t.Fatal(err)
	}
	expected_ok(t,res.Body,"<map[string]string Value>")
	res.Body.Close()
	
}
func Test_params(t *testing.T){
	spk := spark.New()
	spk.Use(spark.ParamsHandler)
	spk.GetRouter().Get("/",func(p spark.Params) string{
		if p["name"] == "rose" {
			return "ok"
		}else {
			return "incorrect"
		}
	})
	ts :=httptest.NewServer(spk)
	defer ts.Close()
	res,err := http.Get(ts.URL+"/?name=rose")
	if err != nil {
		t.Fatal(err)
	}
	expected_ok(t,res.Body,"ok")
}

func expected_ok(t *testing.T,r io.Reader,expect string){
	greeting,err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	x:=len(expect)
	if len(greeting)>=x {
		if (string(greeting[:x]) != expect){
			t.Fatalf("%s",greeting)
		}
	} else {
		t.Fatalf("expected '%s' but '%s'",expect,greeting)
	}
}

