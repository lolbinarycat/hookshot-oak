// +build !

package player

import (
	"github.com/lolbinarycat/hookshot-oak/player/condition"
	"github.com/oakmound/oak/v2/dlog"
)

func (s *PlayerState) DoMap(p *Player) (stateChanged bool) {
	stateChanged = false
	for cond, fun := range s.Map {
		if cond.True(p) {
			nextState := fun.(func(*Player) *PlayerState)(p)
			if nextState != nil {
				dlog.Verb("cond: not nil")
				p.SetState(*nextState)
				stateChanged = true
				break
			}
			dlog.Verb("cond is true")
		}
	}
	return stateChanged
}

func (s *PlayerState) InitMap() {
	for cond := range s.Map {
		switch cond.(type) {
		case condition.Initalized:
			cond.(condition.Initalized).Init()
		}
	}
}
