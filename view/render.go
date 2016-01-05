
package view
import (
//	"fmt"
	"io"
	"os"
	"html/template"
	"path/filepath"
	"regexp"
	"bytes"

	"log"
)


type (
	Config struct{
		// template directory root
		TemplateRoot string
		// not found template name
		NotFoundFile string
		// error template name
		ErrorFile string
	}
	// before run before() before render, data will put
	Render func(w io.Writer,before func()error ,data ...interface{}) error
	RenderMaker interface{
		// return Render by name
		ByName(name string)Render
	}

	renderMakerImpl struct{
		tmpl *template.Template
	}
)

func  DefaultRenderMaker(cfg *Config)(RenderMaker,error){
	root:="html"
	notfoundName:="/notfound\\.(html|tmpl)$"
	errorName:="/error\\.(html|tmpl)$"
	if cfg != nil {

		if len(cfg.TemplateRoot)>0 {
			root = cfg.TemplateRoot
		}
		if len(cfg.NotFoundFile)>0 {
			notfoundName=cfg.NotFoundFile
		}
		if len(cfg.ErrorFile) > 0 {
			errorName=cfg.ErrorFile
		}

	}
	root_len := len(root)

	ext := regexp.MustCompile("\\.(tmpl|html)$")
	nfCompiler := regexp.MustCompile(root+"/"+notfoundName+"$")
	erCompiler :=regexp.MustCompile(root+"/"+errorName+"$")

	nmatch := -1
	ematch := -1
	var b bytes.Buffer

	tmpl := template.New("render234952934___") // create first template
	err := filepath.Walk(root,func(path string,info os.FileInfo,er error)error{

		if er != nil {
			return er
		}
		if info.IsDir() { return nil }
		if ext.MatchString(path){
			if nfCompiler.MatchString(path) {
				if nmatch < 0 {
					nmatch = 1
					notfoundName=path
					return nil
				}
				log.Println("file '%s' match notfound template",path)
			}
			if erCompiler.MatchString(path) {
				if ematch < 0 {
					ematch = 1
					errorName=path
					return nil
				}
				log.Println("file '%s' match error template",path)
			}
			if _,err :=loadTmpl(tmpl,path[root_len+1:],path,&b); err != nil {
				log.Print(err)
			}
		}
		return nil
	})
	if err != nil { return nil,err }

	if nmatch>0{
		_,err := loadTmpl(tmpl,notfoundStr,notfoundName,&b)
		if err != nil {
			log.Println(err)
			template.Must(tmpl.New(notfoundStr).Parse(errHtml))
		}
	} else {
		if len(cfg.NotFoundFile)>0{
			log.Printf("%s is not exits or not ext with (.html|.tmpl), use default template",cfg.NotFoundFile)
		}
		template.Must(tmpl.New(notfoundStr).Parse(errHtml))
	}
	
	if ematch>0{
		_,err := loadTmpl(tmpl,errorStr,errorName,&b)
		if err != nil {
			log.Println(err)
			template.Must(tmpl.New(errorStr).Parse(errHtml))
		}
	} else {
		if len(cfg.ErrorFile)>0{
			log.Printf("%s is not exits or not ext with (.html|.tmpl),use defalt template",cfg.ErrorFile)
		}
		template.Must(tmpl.New(errorStr).Parse(errHtml))
	}
	return &renderMakerImpl{tmpl:tmpl},nil
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
func (rmi *renderMakerImpl)ByName(name string) Render {
	t:=rmi.tmpl.Lookup(name)
	if t == nil {
		return nil
	}
	return newRender(t)
}
func newRender(t *template.Template)Render{
	var buf bytes.Buffer
	return func(w io.Writer,before func()error,data ...interface{})error{
		buf.Reset()

		if len(data)>1{
			if err:=t.Execute(&buf,data); err != nil { return err }
		} else if len(data) == 1{
			if err:=t.Execute(&buf,data[0]); err != nil { return err }
		} else {
			if err:=t.Execute(&buf,nil); err != nil { return err }
		}
		if before != nil {
			if err := before(); err != nil { return err }
		}
		if _,err := io.Copy(w,&buf); err != nil { return err }
		return nil
	}
}

func XmlRender(io.Writer,...interface{})error{return nil}
func JsonRender(io.Writer,...interface{})error{return nil}

/*
 * {"type":<error type>,"error":error}
 */
const (
	notfoundStr = "notfound"
	errorStr = "error"

	errHtml=`<html>
<head><title>{{.type}}</title>
<style type="text/css">
html, body {
font-family: "Roboto", sans-serif;
color: #333333;
background-color: #ea5343;
margin: 0px;
}
h1 {
color: #d04526;
background-color: #ffffff;
padding: 20px;
border-bottom: 1px dashed #2b3848;
}
pre {
margin: 20px;
padding: 20px;
border: 2px solid #2b3848;
background-color: #ffffff;
}
</style>
</head><body>
<h1>{{.type}}</h1>
<pre style="font-weight: bold;">{{.error}}</pre>

</body>
</html>`
)
