package zip

import (
	"io"
	"archive/zip"
	"sync"
)

type (
	Zcallback func(io.ReadCloser,string)error
	delayReader struct{
		once sync.Once
		r io.ReadCloser
		opened bool
		openReader func() (io.ReadCloser,error)
	}
)

func (dr *delayReader)Read(b []byte)(n int,err error){
	dr.once.Do(func(){
		dr.r,err=dr.openReader()
		if err != nil {
			return
		}
		dr.opened=true
	})
	if err != nil {
		return
	}
	return dr.r.Read(b)
}
func (dr *delayReader)Close()error{
	if dr.opened {
		return dr.r.Close()
	}
	return nil
}

func DelayReaderEach(r io.ReaderAt,size int64,cb Zcallback) error{
	zr,err := zip.NewReader(r,size)
	if err != nil { return err}
	for _,f := range zr.File{
		err = cb(&delayReader{openReader:newOpenReader(f),opened:false},f.Name)
		if err != nil{
			return err
		}
	}
	return nil
}
// stop if cb return any error 
func Each(r io.ReaderAt,size int64,cb Zcallback) error{
	zr,err := zip.NewReader(r,size)
	if err != nil { return err}
	for _,f := range zr.File{
		if f.FileInfo().IsDir() {continue}
		suf,er := f.Open()
		if er != nil {
			return nil
		}
		err = cb(suf,f.Name)
		if err != nil{
			return err
		}
	}
	return nil
}
func newOpenReader(f *zip.File) func() (io.ReadCloser,error){
	return func() (io.ReadCloser,error){
		return f.Open()
	}
}
