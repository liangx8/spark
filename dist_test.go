package spark_test

import (
	"github.com/liangx8/spark"
	"testing"
	"net/http"
	"net/http/httptest"
)

func Test_dist(t *testing.T){
	spk := spark.New()
	d:=spark.NewDistribute("action")
	spk.Use(spark.ParamsHandler)
	d.Bind("rose",func()string{
		return "rose1"
	})
	d.Bind("jack",func()string{
		return "jack1"
	})
	spk.GetRouter().Get("/",d.Handler)
	ts :=httptest.NewServer(spk)
	defer ts.Close()
	res,err := http.Get(ts.URL+"/?action=rose")
	if err != nil {
		t.Fatal(err)
	}
	expected_ok(t,res.Body,"rose1")
	res,err = http.Get(ts.URL+"/?action=jack")
	if err != nil {
		t.Fatal(err)
	}
	expected_ok(t,res.Body,"jack1")
}
