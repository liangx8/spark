package main

import (
	"fmt"
	
	"github.com/liangx8/spark"
	"github.com/liangx8/spark/session"
	"github.com/liangx8/spark/view"
)

func main(){
	spk:=spark.New()

	d:=spark.NewDistribute("action")
	// nil means default
	spk.Map(view.ViewReturnHandler(&view.Config{NotFoundFile:"notfound.html"}))

	spk.Use(spark.ParamsHandler)
	spk.Use(session.DefaultHandler())
	d.Bind("rose",func(s session.Session)error{
		if s == nil {
			panic("session 不能是空")
		}
		s.Set("name","rose")
		return fmt.Errorf("rose")
	})
	d.Bind("jack",func(s session.Session)*view.View{

		s.Set("name","jack")
		return view.Html("index.html",map[string]interface{}{
			"tilte":"say hi",
			"message":"hello Jack",
		})
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
	}).Get("/x",d.Handler).NotFound(view.NotFound)
	spk.Run()
}
