package utils

import (
	"errors"
)

//type LastTypeIsNotErr error

var ErrTypeOfLastIsNotError = errors.New("utils: type of last variable is not error")

// func CheckAnd Slice checks the last value, then returns a slice
func CheckAndSlice(vals... interface{}) []interface{} {
	if len(vals) > 1 {
		last := vals[len(vals)-1]
		if  err, ok := last.(error); ok && err != nil {
			panic(err)
		}
	} else if len(vals) == 0 {
		panic("not enough values")
	}
	return vals
}

func Check2(val1 interface{},err error) interface{} {
	if err != nil {
		panic(err)
	}
	return val1
}

func EqualsAny(mArg interface{},sArgs... interface{}) bool {
	for _, arg := range sArgs {
		if mArg == arg {
			return true
		}
	}
	return false
}
