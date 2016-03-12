package view

import (
	"io"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"bytes"
	"log"
)

type (
	Render func(io.Writer,interface{}) error
	RenderMaker interface {
		ByName(string) Render
		Json() Render
		Xml() Render
	}
	Config struct {
		TemplateRoot,TemplateError,TemplateNotFound string
	}
	// a JSON and XML renderer implements
	// ByName is not implemented. Do not call
	TextRenderMaker struct{}
	renderMakerTemplate struct{
		RenderMaker
		tmpls *template.Template
	}
)
func (rmt *renderMakerTemplate)ByName(name string) Render{
	t:=rmt.tmpls.Lookup(name)
	if t == nil {
		return nil
	}
	return func(w io.Writer,data interface{})error{
		if err:=t.Execute(w,data); err != nil { return err }
		return nil
	}
}
// Do not call me
func (*TextRenderMaker)ByName(string)Render{
	panic("do not call me")
}


func (*TextRenderMaker)Json() Render{
	return func(w io.Writer,data interface{}) error {return nil}
}
func (*TextRenderMaker)Xml() Render{
	return func(w io.Writer,data interface{}) error {return nil}
}

func DefaultRenderMaker(cfg *Config) (RenderMaker,error){
	root:="html"
	notfoundName:="/notfound\\.(html|tmpl)$"
	errorName:="/error\\.(html|tmpl)$"
	if cfg != nil {

		if len(cfg.TemplateRoot)>0 {
			root = cfg.TemplateRoot
		}
		if len(cfg.TemplateNotFound)>0 {
			notfoundName=cfg.TemplateNotFound
		}
		if len(cfg.TemplateError) > 0 {
			errorName=cfg.TemplateError
		}

	}
	root_len := len(root)

	ext := regexp.MustCompile("\\.(tmpl|html)$")
	nfCompiler := regexp.MustCompile(root+notfoundName)
	erCompiler :=regexp.MustCompile(root+errorName)


	var b bytes.Buffer

	
	tmpl :=template.New("__render__87d8f") // create first template
	err := filepath.Walk(root,func(path string,info os.FileInfo,er error) error{
		if er != nil {
			return er
		}
		if info.IsDir() { return nil }
		if ext.MatchString(path){
			if nfCompiler.MatchString(path) {
				if _,err :=loadTmpl(tmpl,notfoundStr,path,&b); err != nil{
					log.Printf("load notfound template error with %v",err)
					return nil
				}
				log.Printf("notfound template file '%s'",path)
				return nil
			}
			if erCompiler.MatchString(path) {
				if _,err :=loadTmpl(tmpl,errorStr,path,&b); err != nil{
					log.Printf("load error template error with %v",err)
					return nil
				}
				log.Printf("error template file '%s'",path)
				return nil
			}
			if _,err :=loadTmpl(tmpl,path[root_len+1:],path,&b); err != nil {
				log.Print(err)
			}
		}
		return nil
	})
	if err != nil {
		return nil,err
	}
	return &renderMakerTemplate{RenderMaker:&TextRenderMaker{},tmpls:tmpl},nil
}
func loadTmpl(p *template.Template,name,file string,b *bytes.Buffer) (*template.Template,error){
	log.Printf("load template '%s' as name '%s'",file,name)
	b.Reset()
	if f,err := os.Open(file);err == nil{
		if _,err = io.Copy(b,f); err != nil {
			return nil,err
		}
	} else {
		return nil,err
	}
	return p.New(name).Parse(string(b.Bytes()))
}
const (
	notfoundStr = "notfound"
	errorStr = "error"
)
