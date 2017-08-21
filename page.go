package spark

import (
	"golang.org/x/net/context"
	"html/template"
)

type (
	Page struct{
		Render func(ctx context.Context,data interface{}) error
	}
)
func NewPage(tmpl *template.Template) *Page{
	return &Page{Render:func(ctx context.Context,data interface{}) error{
		w,_,err := ReadHttpContext(ctx)
		if err != nil {
			return err
		}
		return tmpl.Execute(w,data)
	}}
}
