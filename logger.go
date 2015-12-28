package spark
import (
	"log"
	"os"
)
type Logger interface{
	Fatalf(string,...interface{})
	Errorf(string,...interface{})
	Warnningf(string,...interface{})
	Infof(string,...interface{})
}
type simpleLogger struct{
	*log.Logger
}
/*
func (sl *simpleLogger)Fatalf(fmt string,v ...interface{}){
	sl.Fatalf(fmt,v...)
}*/
func (sl *simpleLogger)Errorf(fmt string,v ...interface{}){
	sl.Printf(fmt,v ...)
}
func (sl *simpleLogger)Warnningf(fmt string,v ...interface{}){
	sl.Printf(fmt,v ...)
}
func (sl *simpleLogger)Infof(fmt string,v ...interface{}){
	sl.Printf(fmt,v ...)
}
func DefaultLogger() Logger{
	return &simpleLogger{log.New(os.Stdout,"[spark]",log.Llongfile)}
}
