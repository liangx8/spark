package spark
import (
	"log"
	"os"
	"net/http"
	"time"
)
type Logger interface{
	Fatalf(string,...interface{})
	Errorf(string,...interface{})
	Warnningf(string,...interface{})
	Infof(string,...interface{})
}
type simpleLogger struct{
	*log.Logger
}
/*
func (sl *simpleLogger)Fatalf(fmt string,v ...interface{}){
	sl.Fatalf(fmt,v...)
}*/
func (sl *simpleLogger)Errorf(fmt string,v ...interface{}){
	sl.Printf(fmt,v ...)
}
func (sl *simpleLogger)Warnningf(fmt string,v ...interface{}){
	sl.Printf(fmt,v ...)
}
func (sl *simpleLogger)Infof(fmt string,v ...interface{}){
	sl.Printf(fmt,v ...)
}
func DefaultLogger() Logger{
	return &simpleLogger{log.New(os.Stdout,"[spark]",log.LstdFlags)}
}


func DefaultLogHandler(l Logger,w http.ResponseWriter,req *http.Request,c Context){
	start := time.Now()

	addr := req.Header.Get("X-Real-IP")
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = req.RemoteAddr
		}
	}

	l.Infof("Started %s %s for %s", req.Method, req.URL.Path, addr)

	rw := w.(*ResponseWriter)
	c.Next()

	l.Infof("Completed %v %s in %v\n", rw.Status, http.StatusText(rw.Status), time.Since(start))
}
