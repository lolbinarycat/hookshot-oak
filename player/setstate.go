package player

import (
	"time"
)

type SetStateOptArgs struct {
	SkipEnd bool
}

func (p *Player) SetStateAdv(state PlayerState, opt SetStateOptArgs) {
	if opt.SkipEnd == false && p.State.End != nil {
		p.State.End(p)
	}
	p.StateStartTime = time.Now()
	p.FramesInState = 0

	p.State = state
	//p.State.InitMap()
	if p.State.Start != nil {
		p.State.Start(p)
	}
}

func (p *Player) SetState(state PlayerState) {
	p.SetStateAdv(state, SetStateOptArgs{})
}
