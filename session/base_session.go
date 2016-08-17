package session

import (
	"net/http"
	"sync"
	"reflect"
)
type (
	baseSessionMaker struct{
		w http.ResponseWriter
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
	cookie:= &http.Cookie{Name:SESSION_COOKIE_NAME,Value:id}
	http.SetCookie(sm.w,cookie)
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
func BaseSessionInit(){
	pool :=make(chan map[string]*baseSession,1)
	pool<- make(map[string]*baseSession)
	BuildMaker = func(w http.ResponseWriter,r *http.Request) SessionMaker{
		return &baseSessionMaker{w:w,r:r,sessionPool:pool}
	}
}


