package session

import (
	"net/http"

	"github.com/liangx8/spark"
)

func DefaultHandler()spark.Handler{
	var sc sessionContext
	sc.data=make(map[string]*sessionImpl)
	return func(c spark.Context,w http.ResponseWriter,r *http.Request,l spark.Logger){
		sessionCookie,err := r.Cookie(sessionReqName)

		if err == http.ErrNoCookie {
			sessionCookie = &http.Cookie{
				Name:sessionReqName,
				Value:"",
				Path:"/",
			}
		}
		s:=sc.getOrNew(sessionCookie.Value,l)
		sessionCookie.Value=s.Id()
		s.Lock()
		defer s.Unlock()

		http.SetCookie(w,sessionCookie)
		c.MapTo(s,(*Session)(nil))
		c.Next()
//		s:=sc.getOrNew
	}
}

const (
	sessionReqName = "_gsessionid"
)
