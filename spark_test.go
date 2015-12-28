package spark_test

import (
	"io"
	"io/ioutil"
	"testing"
	"net/http"

	"net/http/httptest"
	"github.com/liangx8/spark"
)


func Test_spark(t *testing.T){
	spk := spark.New()
	spk.Use(func(w http.ResponseWriter){
		panic("painc")
	})
//	spk.GetRouter().Get("/",func(w http.ResponseWriter){
//		fmt.Fprint(w,"ok")
//	})
	ts :=httptest.NewServer(spk)
	defer ts.Close()
	res,err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	expected_ok(t,res.Body)
}

func expected_ok(t *testing.T,r io.Reader){
	greeting,err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(greeting)>=2 {
		if (string(greeting[:2]) != "ok"){
			t.Fatalf("%s",greeting)
		}
	}
}

