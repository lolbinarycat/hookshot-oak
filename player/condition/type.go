package condition

import (
	"github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/dlog"
)

type Condition interface {
	True() bool
}

type Initalized interface {
	Init()
}

// KeyDown is a condition that is true if a key is held down
type KeyDown string

func (k KeyDown) True() bool {
	return oak.IsDown(string(k))
}

type FramesElapsed struct {
	N, c int
}

func (f *FramesElapsed) Init() {
	f.c = 0
}

func (f *FramesElapsed) True() bool {
	dlog.Verb("cond:",f.c,f.N)
	f.c++
	return f.c >= f.N
}

type Func func () bool

func (f Func) True() bool {
	return f()
}
