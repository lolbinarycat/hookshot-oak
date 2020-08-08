package player

import (
	"time"
)

type SetStateOptArgs struct {
	SkipEnd bool
}

func (p *Player) SetStateAdv(state interface{}, opt SetStateOptArgs) {
	var st *PlayerState
	switch state.(type) {
	case *PlayerState:
		st = state.(*PlayerState)
	case PlayerState:
		tmp := state.(PlayerState)
		st = &tmp
	default:
		panic("bad type as argument to SetState")
	}
	if opt.SkipEnd == false && p.State.End != nil {
		p.State.End(p)
	}
	p.StateStartTime = time.Now()
	p.FramesInState = 0

	p.State = st
	//p.State.InitMap()
	if p.State.Start != nil {
		p.State.Start(p)
	}
}

func (p *Player) SetState(state interface{}) {
	p.SetStateAdv(state, SetStateOptArgs{})
}
