package view

import (
	"io"
)

type (
	Render func(io.Writer,interface{}) error
	RenderMaker interface {
		ByName(string) Render
		Json() Render
		Xml() Render
	}
	Config struct {
		TemplateDir,TemplateError,TemplateNotFound string
	}
)
