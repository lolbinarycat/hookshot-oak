package player

import (
	"time"
)

type SetStateOptArgs struct {
	SkipEnd bool
}

func (p *Player) SetStateAdv(state PlayerState, opt SetStateOptArgs) {
	if opt.SkipEnd == false {
		p.State.End(p)
	}
	p.StateStartTime = time.Now()

	p.State = state
	p.State.InitMap()
	p.State.Start(p)
}

func (p *Player) SetState(state PlayerState) {
	p.SetStateAdv(state, SetStateOptArgs{})
}
