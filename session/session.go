package session

import (
	"net/http"
	"errors"
	"golang.org/x/net/context"


	"github.com/liangx8/spark"
)

type (

	Session interface{
		Get(key string,ptr interface{}) bool
		Put(key string,ptr interface{})
		Id() string
	}
	key int
	SessionMaker func(context.Context,string)Session

	

)
const (
	sessionKey key = iota
)
// Get Session
func GetSession(ctx context.Context)(Session,error){
	s,ok:=ctx.Value(sessionKey).(Session)
	if ok {
		return s,nil
	}
	return nil,ErrNotSessionContext
}
// Usage:
//   spk:=New(func(_ *http.Request)context.Context{return context.Background})
//   spk.AddChain(CreateSessionChain(SimpleSessionMaker))
func CreateSessionChain(maker SessionMaker) func(context.Context,spark.HandleFunc){
	return func(ctx context.Context,chain spark.HandleFunc){

        w,r,err:=spark.ReadHttpContext(ctx)
		if err != nil {
			panic(err)
		}
		id := obtainId(r)
		s:=maker(ctx,id)
		if s == nil {
			panic("Session is nil,perhps SessionMaker is not a complete implements")
		}
			
		setCookie(w,s.Id())
		chain(context.WithValue(ctx,sessionKey,s))
	}
}

func obtainId(req *http.Request)string{
	cookie,err:=req.Cookie(SESSION_COOKIE_NAME)

	if err == http.ErrNoCookie {
		return ""
	}

	return cookie.Value
}
func setCookie(w http.ResponseWriter,id string){
	cookie:= &http.Cookie{Name:SESSION_COOKIE_NAME,Value:id}
	http.SetCookie(w,cookie)
}

var ErrNotSessionContext = errors.New("Not session context")

const (
	SESSION_COOKIE_NAME = "_pfa_SESSION_ID"
)


