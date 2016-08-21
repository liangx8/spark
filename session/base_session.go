package session

import (

	"sync"
	"reflect"


	"golang.org/x/net/context"

)

type (
	baseSession struct{
		m     sync.Mutex
		pool  map[string]reflect.Value
		getId func() string
	}
)


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
// the SimpleSessionMaker create a never expire Session
func SimpleSessionMaker(_ context.Context,id string)Session{
	var sid string
	mutex.Lock()
	defer mutex.Unlock()
	var s *baseSession
	ok := false
	if id != ""{
		s,ok = pool[id]
	}
	if !ok {
		sid = UniqueId()
		s = &baseSession{
			pool:make(map[string]reflect.Value),
			getId:func()string{ return sid },
		}
		pool[sid]=s
	}
	return s

}

var (
	pool map[string]*baseSession = make(map[string]*baseSession)
	mutex sync.Mutex
)
