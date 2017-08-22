package spark

import (
	"golang.org/x/net/context"
	"html/template"
	"io"
	"errors"
)

type (
	View struct{
		render func(w io.Writer,data interface{})error
	}
)
func NewTemplateView(tmpl *template.Template) *View{
	return &View{render:func(w io.Writer ,data interface{}) error{
		return tmpl.Execute(w,data)
	}}
}
func (p View)Render(ctx context.Context, data interface{}) error{
		w,_,err := ReadHttpContext(ctx)
		if err != nil {
			return err
		}
	return p.render(w,data)
}
func NewStringView() *View{
	return &View{render:func(w io.Writer, data interface{})error{
		value,ok := data.(string)
		if ok {
			_,err:=w.Write([]byte(value))
			if err != nil {
				return err
			}
			return nil
		}
		return NotStringError
	}}
}

var (
	NotStringError = errors.New("Not string")
)
