package session

import (
	"net/http"
)

type (
	Session interface{
		Get(key string,ptr interface{}) bool
		Put(key string,ptr interface{})
		Id() string
	}
	SessionMaker interface{
		New() Session
		Get(string) Session
		//is session id valid ?
		IsValid(string) bool
	}
)

var BuildMaker func(http.ResponseWriter,*http.Request) SessionMaker
func Get(w http.ResponseWriter,req *http.Request) Session{
	maker := BuildMaker(w,req)
	id := obtainId(req)
	if id == "" {
		return maker.New()
	}
	s := maker.Get(id)
	if s != nil {
		return s
	}
	return maker.New()
}
func obtainId(req *http.Request)string{

	cookie,err:=req.Cookie(SESSION_COOKIE_NAME)
	if err == http.ErrNoCookie {
		return ""
	}
	return cookie.Value
}
const (
	SESSION_COOKIE_NAME = "_pfa_SESSION_ID"
)

