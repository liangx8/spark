package spark

import (
	"fmt"
)

func ErrorWrap(previous error,msg string) error{
	return fmt.Errorf("%s wrapped error:\n%v",msg,previous)
}
