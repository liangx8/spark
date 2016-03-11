package view_test

import (
	"testing"
	"os"
	"fmt"
	"bytes"

	"github.com/liangx8/spark/view"

)
func Test_defaultRenderMaker(t *testing.T){
	_,err:=view.DefaultRenderMaker(nil)
	if err != nil {
		t.Fatalf("not expected %v",err)
	}
}
func Test_render(t *testing.T){
	var buf bytes.Buffer
	defer cleanDir()
	if err:=prepareDir();err != nil {
		t.Fatal(err)
	}
	maker,err:=view.DefaultRenderMaker(&view.Config{TemplateRoot:test_dir})
	if err != nil {
		t.Fatal(err)
	}
	for i:=0;i< 5;i ++ {
		render := mustGetRender(t,maker,fmt.Sprintf("test%d.tmpl",i))
		mustRender(t,render,&buf,nil)
		expectedContent(t,&buf,fmt.Sprintf("test%d content",i))
	}
	render := mustGetRender(t,maker,"test6.tmpl")
	mustRender(t,render,&buf,1024)
	expectedContent(t,&buf,"test6 1024")
	render = mustGetRender(t,maker,"error")
	mustRender(t,render,&buf,nil)
	expectedContent(t,&buf,"error")
	render = mustGetRender(t,maker,"notfound")
	mustRender(t,render,&buf,nil)
	expectedContent(t,&buf,"notfound")
}
func mustRender(t *testing.T,r view.Render,buf *bytes.Buffer,data interface{}){
	buf.Reset()
	if err := r(buf,data);err != nil {
		t.Fatal(err)
	}
	
}
func mustGetRender(t *testing.T,maker view.RenderMaker,name string) view.Render{
	render:=maker.ByName(name)
	if render == nil {
		t.Fatalf("expected %s",name)
	}
	return render
}
func expectedContent(t *testing.T, buf *bytes.Buffer,content string){
	s := string(buf.Bytes())
	if s != content {
		t.Fatalf("expected '%s', but '%s'",content,s)
	}
}
func cleanDir(){
	for i:=0 ; i< 5; i++ {
		file := fmt.Sprintf("%s/test%d.tmpl",test_dir,i)
		os.Remove(file)

	}
	os.Remove(fmt.Sprintf("%s/test6.tmpl",test_dir))
	os.Remove(fmt.Sprintf("%s/error.tmpl",test_dir))
	os.Remove(fmt.Sprintf("%s/notfound.html",test_dir))

	os.Remove(fmt.Sprintf("%s/sub/test0.tmpl",test_dir))
	os.Remove(test_dir+"/sub")
	os.Remove(test_dir)
}
func prepareDir() error{

	if err:=os.Mkdir(test_dir,os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	for i:=0 ; i< 5; i++ {
		file := fmt.Sprintf("%s/test%d.tmpl",test_dir,i)
		f,err := os.Create(file)
		defer f.Close()
		if err != nil {
			return err
		}
		fmt.Fprintf(f,"test%d content",i)
	}
	{
		f,err := os.Create(fmt.Sprintf("%s/test6.tmpl",test_dir))
		if err != nil {
			return err
		}
		defer f.Close()
		fmt.Fprintf(f,"test6 {{.}}")
	}
	{
		f,err := os.Create(fmt.Sprintf("%s/error.tmpl",test_dir))
		if err != nil { return err }
		defer f.Close()
		fmt.Fprintf(f,"error")
	}
	{
		f,err := os.Create(fmt.Sprintf("%s/notfound.html",test_dir))
		if err != nil { return err }
		defer f.Close()
		fmt.Fprintf(f,"notfound")
	}
	{
		if err:=os.Mkdir(fmt.Sprintf("%s/sub",test_dir),os.ModeDir|os.ModePerm); err != nil {return err}
		f,err := os.Create(fmt.Sprintf("%s/sub/test0.tmpl",test_dir))
		if err != nil {return err}
		defer f.Close()
		fmt.Fprintf(f,"sub/test0")
	}
	return nil
}


const (
	test_dir = "/tmp/test_templates"
)
