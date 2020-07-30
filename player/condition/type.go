package condition

import (
	"github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/dlog"
)

// we use an empty interface instead of *Player to prevent circular imports 
type any = interface{}

type Condition interface {
	True(interface{}) bool
}

type Initalized interface {
	Init()
}

// KeyDown is a condition that is true if a key is held down
type KeyDown string

func (k KeyDown) True(any) bool {
	return oak.IsDown(string(k))
}

type FramesElapsed struct {
	N, c int
}

func (f *FramesElapsed) Init() {
	f.c = 0
}

func (f *FramesElapsed) True(any) bool {
	dlog.Verb("cond:",f.c,f.N)
	f.c++
	return f.c >= f.N
}

type Func func (any) bool

func (f Func) True(a any) bool {
	return f(a)
}
