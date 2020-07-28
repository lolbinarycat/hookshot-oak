package condition

import (
	"github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/event"
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
	Target, current int
}

func (f FramesElapsed) Init() {
	event.GlobalBind(func (_ int, _ interface{}) int {
		f.current++
		if f.current > f.Target {
			return event.UnbindSingle
		}
		return 0
	}, event.Enter)
}

func (f FramesElapsed) True() bool {
	return f.current >= f.Target
}

type Func func () bool

func (f Func) True() bool {
	return f()
}
