package view
import (
	"fmt"
	"io"
	"os"
	"html/template"
	"path/filepath"
	"regexp"
)
type RenderMaker interface{
	ByName( string) (Render,error)
}
func DefaultRenderMaker(root string) (RenderMaker,error){
	files := make([]string,0,10)
	ext := regexp.MustCompile("\\.(tmpl|html)$")
	var tmpl *template.Template
	err := filepath.Walk(root,func(path string,info os.FileInfo,er error)error{
		if er != nil {
			return er
		}
		if info.IsDir() { return nil }
		if ext.MatchString(path){
			files = append(files,path)
		}
		return nil
	})
	if err != nil { return nil,err }
	tmpl,err = template.ParseFiles(files...)
	if err != nil { return nil,err }
	return &renderMakerImpl{tmpl},nil
}
type Render func(io.Writer,...interface{})

type renderMakerImpl struct{
	tmpl *template.Template
}
func (rmi *renderMakerImpl)ByName(name string) (Render,error) {
	t:=rmi.tmpl.Lookup(name)
	if t == nil {
		return nil,fmt.Errorf("Don't know how to make Render '%s'",name)
	}
	return func(w io.Writer,data ...interface{}){
		if len(data)>1{
			t.Execute(w,data)
		} else if len(data) == 1{
			t.Execute(w,data[0])
		} else {
			t.Execute(w,nil)
		}
	},nil
}
func XmlRender(io.Writer,...interface{}){}
func JsonRender(io.Writer,...interface{}){}
