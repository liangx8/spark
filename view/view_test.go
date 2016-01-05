package view_test
import (
	"github.com/liangx8/spark/view"
	"testing"

	"os"
	"fmt"
	"bytes"
	"regexp"
)
/*

type Config struct{
	// template directory root
	TemplateRoot string
	// not found template name
	NotFoundName string
	// error template name
	ErrorName string
}
*/
func Test_re(t *testing.T){
	prefix := "aa/bb/cc"
	error_c := regexp.MustCompile("/error\\.(html|tmpl)$")
	like := make([]string,0,10)
	strs := []string{
		prefix+"/ee.html",
		prefix+"/error.html",
		prefix+"/xerror.html",
		prefix+"/xx.html",
		prefix+"/xxx/error.html",
	}
	match := -1
	for i,x := range strs {
		if error_c.MatchString(x) {
			if match < 0 {
				match = i
				t.Logf("pack '%s'",x)
			} else {
				like = append(like,x)
			}
		}
	}
	if len(like)>0 {
		t.Log(like)
	}
	if match != 1 {
		t.Fatalf("match at %d expected 1",match)
	}
	error_c = regexp.MustCompile(prefix+"/error.html$")
	match = -1
	like = like[0:0]
	for i,x := range strs {
		if error_c.MatchString(x) {
			if match < 0 {
				match = i
				t.Logf("pack '%s'",x)
			} else {
				like = append(like,x)
			}
		}
	}
	if len(like)>0 {
		t.Log(like)
	}
	if match != 1 {
		t.Fatalf("match at %d expected 1",match)
	}
}
func Test_render(t *testing.T){
	defer cleanDir()
	var b bytes.Buffer
	var r view.Render
	if err:=prepareDir();err != nil {
		t.Log(os.Getwd())
		t.Fatal(err)
	}
/*
	maker,err:=view.DefaultRenderMaker(&view.Config{
		TemplateRoot:test_dir,
		NotFoundName:"nofound.html",
		ErrorName:"error.tmpl",
	})
*/
	maker,err:=view.DefaultRenderMaker(&view.Config{
		TemplateRoot:test_dir,
		NotFoundFile:"xxx.html",
	})
	if err != nil {
		t.Fatal(err)
	}
	if r =maker.ByName("test0.tmpl"); r == nil { t.Fatal("unexpected nil")}
	r(&b,nil)
	expectedString(t,"test0 content",string(b.Bytes()))
	if r =maker.ByName("test1.tmpl"); r == nil { t.Fatal("unexpected nil")}
	b.Reset()
	r(&b,nil)
	expectedString(t,"test1 content",string(b.Bytes()))
	if r =maker.ByName("test6.tmpl"); r == nil { t.Fatal("unexpected nil")}
	b.Reset()
	r(&b,nil,100)
	expectedString(t,"test6 100",string(b.Bytes()))
	if r =maker.ByName("error.tmpl"); r == nil { t.Fatal("unexpected nil")}
	b.Reset()
	r(&b,nil)
	expectedString(t,"error",string(b.Bytes()))
	if r =maker.ByName("error"); r == nil { t.Fatal("error unexpected nil")}
	b.Reset()
	r(&b,nil)
	if string(b.Bytes()) == "error"  {t.Fatal("error expected a default value.")}
	if r =maker.ByName("notfound"); r == nil { t.Fatal("notfound unexpected nil")}

	if r =maker.ByName("sub/test0.tmpl"); r == nil { t.Fatal("unexpected nil")}
	b.Reset()
	r(&b,nil)
	expectedString(t,"sub/test0",string(b.Bytes()))

	maker,err =view.DefaultRenderMaker(&view.Config{
		TemplateRoot:test_dir,
		NotFoundFile:"notfound.html",
		ErrorFile:"error.tmpl",
	})

	if r =maker.ByName("error"); r == nil { t.Fatal("error unexpected nil")}
	b.Reset()
	r(&b,nil)
	expectedString(t,"error",string(b.Bytes()))

	if r =maker.ByName("notfound"); r == nil { t.Fatal("notfound unexpected nil")}
	b.Reset()
	r(&b,nil)
	expectedString(t,"notfound",string(b.Bytes()))
}

func expectedString(t *testing.T, s1,s2 string){
	if s1 != s2 {
		t.Fatalf("expected %s, but %s",s1,s2)
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
	test_dir = "test_templates"
)
