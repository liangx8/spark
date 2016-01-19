package spark
import (
	"log"
//	"os"
	"net/http"
	"time"
)

func DefaultLogHandler(l *log.Logger,w http.ResponseWriter,req *http.Request,c Context){
	start := time.Now()

	addr := req.Header.Get("X-Real-IP")
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = req.RemoteAddr
		}
	}

	l.Printf("Started %s %s for %s", req.Method, req.URL.Path, addr)

	rw := w.(*responseWriter)
	c.Next()

	l.Printf("Completed %d %s in %v\n", rw.Status, http.StatusText(rw.Status), time.Since(start))
}
