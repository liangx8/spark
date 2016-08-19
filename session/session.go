package session

import (

	"errors"
	"golang.org/x/net/context"

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
	key int

)
const (
	sessionKey key = iota
)
func GetSession(ctx context.Context)(Session,error){
	s,ok:=ctx.Value(sessionKey).(Session)
	if ok {
		return s,nil
	}
	return nil,ErrNotSessionContext
}

var ErrNotSessionContext = errors.New("Not session context")
