package spark
import (
	"log"
	"os"
	"net/http"
	"time"
)
type (
	LogAdaptor func(LogLevel,string,...interface{})
	LogLevel int
)
const (
	INFO LogLevel= iota
	WARNNING
	ERROR
	FATAL
)
func DefaultLogHandler(l LogAdaptor,w http.ResponseWriter,req *http.Request,c Context){
	start := time.Now()

	addr := req.Header.Get("X-Real-IP")
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = req.RemoteAddr
		}
	}

	l(INFO,"Started %s %s for %s", req.Method, req.URL.Path, addr)

	rw := w.(*responseWriter)
	c.Next()

	l(INFO,"Completed %d %s in %v\n", rw.Status, http.StatusText(rw.Status), time.Since(start))
}

func DefaultLog() LogAdaptor{
	l := log.New(os.Stdout,"[spark] ",log.LstdFlags)
	return func(lvl LogLevel,fmt string,args ...interface{}){
		switch lvl{
		case INFO,WARNNING:
			l.Printf(fmt,args...)
		case ERROR:
			l.Fatalf(fmt,args...)
		case FATAL:
			l.Panicf(fmt,args...)
		}
	}
}
var EmptyLogAdaptor LogAdaptor = func(level LogLevel,fmt string, objs ...interface{}){
	// do nothing
}
