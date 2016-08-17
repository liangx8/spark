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

var BuildMaker func(*http.Request) SessionMaker
func Get(w http.ResponseWriter,req *http.Request) Session{
	maker := BuildMaker(req)
	id := obtainId(req)
	if id == "" {
		s := maker.New()
		setCookie(w,s.Id())
		return s
	}
	s := maker.Get(id)
	if s != nil {
		return s
	}
	s=maker.New()
	setCookie(w,s.Id())
	return s
}
func setCookie(w http.ResponseWriter,id string){
	cookie:= &http.Cookie{Name:SESSION_COOKIE_NAME,Value:id}
	http.SetCookie(w,cookie)
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

