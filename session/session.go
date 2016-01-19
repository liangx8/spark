package session
import (
	"time"
	"sync"
	"errors"
	"log"
)
var SessionExpire time.Duration
type (
	sessionContext map[string]*sessionImpl

	Session interface{
//		Lock()
//		Unlock()
		// return value for key
		Get(key string) interface{}
		// set value
		Set(key string, value interface{})
		// return Id
		Id() string
		// Invalidate Session
		Invalidate()
	}

	sessionImpl struct {
		sync.Mutex
		data map[string]interface{}
		expire *time.Timer
		removeFromContext func()
		returnId func()string
		//	valid bool
	}
)

func (sc *sessionContext)getOrNew(id string,l *log.Logger)*sessionImpl{
	s,ok := (*sc)[id]
	if ok {
		if s.expire.Reset(SessionExpire) {
			return s
		}
	}
	id = newSessionId()
	s = &sessionImpl{
		returnId:func()string{return id},
		data:make(map[string]interface{}),
		expire:time.AfterFunc(SessionExpire,func(){
			delete(*sc,id)
			log.Printf("session %s expires",id)	
		}),
		removeFromContext:func(){
			if s.expire.Stop() {
				delete(*sc,id)
				l.Printf("session %s invalidate",id)
			}
		},
	}
	(*sc)[id]=s
	l.Printf("Create a new session %s",id)	
	return s
}

func (s *sessionImpl)Invalidate(){

	s.removeFromContext()
}
func (s *sessionImpl)Get(key string) interface{}{

	value := s.data[key]
	return value
}

func (s *sessionImpl)Set(key string,value interface{}){
	s.data[key]=value

}

func (s *sessionImpl)Id()string{
	return s.returnId()
}

//func (s *sessionImpl)

func newSessionId() string{
	i := time.Now().UnixNano()
	buf := make([]byte,0,30)

	v := i % (26 + 26 + 10)
	for i > 0 {
		d := (i+v) % (26 + 26 + 10)
		i = i / (26 + 26 + 10)
		if d < 10 {
			buf = append(buf,byte(d + 0x30))
		} else {
			d = d - 10
			if d < 26 {
				buf = append(buf,byte(d + 0x41))
			} else {
				d = d- 26
				buf = append(buf,byte(d + 0x61))
			}
		}
	}
	return string(buf)
}


var ErrSessionInvalid = errors.New("Session had been invalided")
