package spark
import (
	"reflect"
)

type ReturnHandler func(c Context,retval []reflect.Value)
