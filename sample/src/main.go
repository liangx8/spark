package main

import (
	"github.com/liangx8/spark"
	"github.com/liangx8/spark/session"
)

func main(){
	spk:=spark.New()

	d:=spark.NewDistribute("action")
	spk.Use(spark.ParamsHandler)
	spk.Use(session.DefaultHandler())
	d.Bind("rose",func(s session.Session)string{
		name,_ :=s.Get("name")
		s.Set("name","rose")
		return name.(string) + " rose"
	})
	d.Bind("jack",func(s session.Session)string{
		name,_ :=s.Get("name")
		s.Set("name","jack")
		return name.(string) + " jack"
	})

	spk.GetRouter().Get("/",
		func()string{
			return "hello world\n"
		},
		func()string{
			return "nice to meet you"
		},
	).Get("/panic",func(){
		panic("raise a panic")
	}).Get("/x",d.Handler)
	spk.Run()
}
