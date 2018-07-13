package spark

import (
	"golang.org/x/net/context"
	"html/template"
	"io"
	"errors"
	"encoding/json"
	"gopkg.in/yaml.v2"
)

type (
	view struct{
		render func(w io.Writer,data interface{})error
	}
	View interface{
		Render(ctx context.Context,data interface{}) error
	}
)
func NewTemplateView(tmpl *template.Template) View{
	return &view{render:func(w io.Writer ,data interface{}) error{
		return tmpl.Execute(w,data)
	}}
}
func (p *view)Render(ctx context.Context, data interface{}) error{
		w,_,err := ReadHttpContext(ctx)
		if err != nil {
			return err
		}
	return p.render(w,data)
}
func NewStringView() View{
	return &view{render:func(w io.Writer, data interface{})error{
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
func NewJSONView() View{
	return &view{
		render:func(w io.Writer,data interface{})error{
			en:=json.NewEncoder(w)
			return en.Encode(data)
		},
	}
}
func NewYAMLView() View{
	return &view{
		render:func(w io.Writer,data interface{}) error{
			en:=yaml.NewEncoder(w)
			defer en.Close()
			return en.Encode(data)
		},
	}
}
var (
	NotStringError = errors.New("Not string")
)
