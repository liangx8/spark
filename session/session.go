package session
import (
	"time"
	"sync"
	"errors"
	"fmt"

	"github.com/liangx8/spark"
)

type sessionContext struct{
	data map[string]*sessionImpl
}

func (sc *sessionContext)getOrNew(id string,log spark.Logger)Session{
	s,ok := sc.data[id]
	if ok {
		return s
	}
	s = &sessionImpl{
		id:fmt.Sprint(time.Now().UnixNano()),
		data:make(map[string]interface{}),
		expire:time.NewTimer(15 * time.Minute),
	}
	sc.data[s.id]=s
	log.Infof("Create a new session %s",s.id)	
	return s
}

type Session interface{
	Lock()
	Unlock()
	// return value for key or return ErrSessionInvalid
	Get(key string) (interface{},error)
	// set value or return ErrSessionInvalid
	Set(key string, value interface{}) error
	// return Id
	Id() string
	// Invalidate Session
	Invalidate()
}

type sessionImpl struct {
	sync.Mutex
	id string
	data map[string]interface{}
	expire *time.Timer
//	valid bool
}

func (s *sessionImpl)Invalidate(){
	s.expire.Stop()
}
func (s *sessionImpl)Get(key string) (interface{},error){

	value := s.data[key]
	if !s.expire.Reset(15 * time.Minute){
		return nil,ErrSessionInvalid
	}
	return value,nil
}

func (s *sessionImpl)Set(key string,value interface{})error{

	s.data[key]=value
	if !s.expire.Reset(15 * time.Minute){
		return ErrSessionInvalid
	}
	return nil
}

func (s *sessionImpl)Id()string{
	return s.id
}

//func (s *sessionImpl)


var ErrSessionInvalid = errors.New("Session had been invalided")
