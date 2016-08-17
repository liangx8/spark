package session

import (
	"crypto/rand"
)

func UniqueId() string{
	buf := make([]byte,20)
	_,err:=rand.Reader.Read(buf)
	if err != nil {
		panic(err)
	}
	for i := range buf{
		buf[i]=table[buf[i] % 62]
	}
	return string(buf)
}

var table []byte=[]byte("01234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
