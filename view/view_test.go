package view_test
import (
	"github.com/liangx8/spark/view"
	"testing"

	"os"
	"fmt"
	"bytes"
)

func Test_render(t *testing.T){
	defer cleanDir()
	var b bytes.Buffer
	var r view.Render
	if err:=prepareDir();err != nil {
		t.Log(os.Getwd())
		t.Fatal(err)
	}
	
	maker,err:=view.DefaultRenderMaker(test_dir)
	if err != nil {
		t.Fatal(err)
	}
	if r,err =maker.ByName("test0.tmpl"); err != nil { t.Fatal(err)}
	r(&b)
	expectedString(t,"test0 content",string(b.Bytes()))
	if r,err =maker.ByName("test1.tmpl"); err != nil { t.Fatal(err)}
	b.Reset()
	r(&b)
	expectedString(t,"test1 content",string(b.Bytes()))
	if r,err =maker.ByName("test6.tmpl"); err != nil { t.Fatal(err)}
	b.Reset()
	r(&b,100)
	expectedString(t,"test6 100",string(b.Bytes()))
	t.Fatal(view.DefaultRenderMaker("xx"))

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
	os.Remove(test_dir)
}
func prepareDir() error{

	if err:=os.Mkdir(test_dir,os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	for i:=0 ; i< 5; i++ {
		file := fmt.Sprintf("%s/test%d.tmpl",test_dir,i)
		f,err := os.Create(file)
		if err != nil {
			return err
		}
		fmt.Fprintf(f,"test%d content",i)
		f.Close()
	}
	f,err := os.Create(fmt.Sprintf("%s/test6.tmpl",test_dir))
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintf(f,"test6 {{.}}")

	return nil
}

const (
	test_dir = "test_templates"
)
