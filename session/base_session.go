package session

import (
	"net/http"
	"sync"
	"reflect"

	"golang.org/x/net/context"
	"github.com/liangx8/spark"
)
type (
	baseSessionMaker struct{
		r *http.Request
		sessionPool chan map[string]*baseSession
	}
	baseSession struct{
		m sync.Mutex
		pool map[string]reflect.Value
		getId func() string
	}
)

func (sm *baseSessionMaker)New() Session{
	h := <-sm.sessionPool
	defer func(){sm.sessionPool<- h}()
	id := UniqueId()

	bs:=&baseSession{
		getId:func()string{return id},
		pool:make(map[string]reflect.Value),
	}
	h[id]=bs
	return bs
}
func (sm *baseSessionMaker)Get(id string) Session{
	h := <-sm.sessionPool
	defer func(){sm.sessionPool<- h}()
	s,ok :=h[id]
	if ok {return s}
	return nil
}
func (sm *baseSessionMaker)IsValid(id string) bool{
	return false
}

func (bs *baseSession)Get(key string, ptr interface{}) bool{
	bs.m.Lock()
	defer bs.m.Unlock()
	v,ok := bs.pool[key]
	if !ok {return false}
	ptr_v := reflect.ValueOf(ptr).Elem()
	ptr_v.Set(v)
	return true
}
func (bs *baseSession)Put(key string, ptr interface{}){
	bs.m.Lock()
	defer bs.m.Unlock()
	bs.pool[key]=reflect.ValueOf(ptr).Elem()
}
func (bs *baseSession)Id() string{
	return bs.getId()
}

func CreateSessionChain()func(context.Context,spark.HandleFunc){
	pool :=make(chan map[string]*baseSession,1)
	pool<- make(map[string]*baseSession)
	return func(ctx context.Context, chain spark.HandleFunc){
		w,r,err:=spark.ReadHttpContext(ctx)
		if err != nil {
			panic(err)
		}
		id := obtainId(r)
		h:= <-pool
		defer func(){
			pool<- h
		}()
		s,ok := h[id]
		if !ok {
			id = UniqueId()
			s = &baseSession{
				pool:make(map[string]reflect.Value),
				getId:func()string{ return id },
			}
			h[id]=s
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

const (
	SESSION_COOKIE_NAME = "_pfa_SESSION_ID"
)
