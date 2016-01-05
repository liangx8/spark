package session
import (
	"time"
	"sync"
	"errors"
	"fmt"

	"github.com/liangx8/spark"
)

type sessionContext map[string]*sessionImpl


func (sc *sessionContext)getOrNew(id string,log spark.Logger)Session{
	s,ok := (*sc)[id]
	if ok {
		return s
	}
	id = fmt.Sprint(time.Now().UnixNano())
	s = &sessionImpl{
		returnId:func()string{return id},
		data:make(map[string]interface{}),
		expire:time.AfterFunc(15 * time.Minute,func(){
			delete(*sc,id)
			log.Infof("session %s expires",id)	
		}),
		removeFromContext:func(){
			delete(*sc,id)
			log.Infof("session %s invalidate",id)
		},
	}
	(*sc)[id]=s
	log.Infof("Create a new session %s",id)	
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
	data map[string]interface{}
	expire *time.Timer
	removeFromContext func()
	returnId func()string
//	valid bool
}

func (s *sessionImpl)Invalidate(){
	s.expire.Stop()
	s.removeFromContext()
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
	return s.returnId()
}

//func (s *sessionImpl)


var ErrSessionInvalid = errors.New("Session had been invalided")
