package order

import (
	"sort"
)

type (
	// 把对象按照给于的 seq 顺序提取 
	Order interface{
		// seq 顺序

		// obj 对象
		Add(seq int, obj interface{})
		NewSequence() Sequence
	}
	Sequence interface{
		Next() bool
		Object() interface{}
	}
	order struct{
		sort.IntSlice
		objs map[int][]interface{}
	}
	sequence struct{
		next func() bool
		obj func() interface{}
	}
)

func New() Order{
	return &order{objs:make(map[int][]interface{})}
}

func (o *order)Add(seq int,obj interface{}){
	eobjs,exists := o.objs[seq]
	var objs []interface{}
	if exists {
		objs = append(eobjs,obj)
	} else {
		objs = make([]interface{},1)
		o.IntSlice = append(o.IntSlice,seq)
		objs[0]=obj
	}
	o.objs[seq]=objs
}
func (o *order)NewSequence() Sequence{
	acount := len(o.IntSlice)
	if acount == 0 {
		return emptySequence
	}
	o.Sort()
	aidx :=0
	objs := o.objs[o.IntSlice[0]][:]
	ocount := len(objs)
	oidx := -1
	return &sequence{
		next:func()bool{

			oidx ++

			if oidx >= ocount {
				aidx ++
				oidx = 0
			}

			if aidx >= acount {
				return false
			}

			objs = o.objs[o.IntSlice[aidx]][:]
			ocount=len(objs)

			return true
		},
		obj:func()interface{}{
			if oidx < 0 {
				panic("Beginning of Sequence Error, you should call Next() first")
			}
			return objs[oidx]
		},
	}
}
var emptySequence = &sequence{
	next:func()bool{
		return false
	},
	obj:func()interface{}{
		return nil
	},
}
func (s sequence)Next()bool{
	return s.next()
}
func (s sequence)Object()interface{}{
	return s.obj()
}
