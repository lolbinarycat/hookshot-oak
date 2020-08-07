package player

import (
	"time"
	"fmt"

	"github.com/oakmound/oak/v2/dlog"
)

// denil returns a modified version of a playerstate with nil functions replaced
// with empty functions, preventing segfaults from happening if they are called.
// It is designed to be called on a struct literal when setting a value
func (s PlayerState) denil() PlayerState {
	if s.Start == nil {
		s.Start = func(p *Player) {}
	}
	if s.LLoop == nil {
		s.LLoop = func(p *Player) *PlayerState {return nil}
	}
	if s.Loop == nil {
		s.Loop = func(p *Player) {}
	}
	if s.End == nil {
		s.End = func(p *Player) {}
	}
	if (s.NextState == nil) {
		s.NextState = &AirState
	}
	if (s.MaxDuration == 0) {
		s.MaxDuration = time.Minute * 20
	}
	for k, v := range s.Map {
		switch v.(type) {
		case *PlayerState:
			state := v.(*PlayerState)
			s.Map[k] = func(_ *Player) *PlayerState {
				return state
			}
		case PlayerStateMapFunc, (func(*Player) *PlayerState) :
			break // do nothing
		default:
			panic(fmt.Sprintf("unexpected type %T as value to StateMap",v))
		}
	}
	return s
}



// initStates is called at the start of main().
// this is to stop an initialization error.
func init() {
	AirState = PlayerState{
		Loop:AirStateLoop,
		Start:func(p *Player) {},
		End:func(p *Player) {},
	}.denil()
	GroundState = PlayerState{
		Loop: GroundStateLoop,
	}.denil()
	ItemCarryGroundState = PlayerState{
		Start: func(p *Player) {
			if p.HeldObj == nil {
				dlog.Error("HeldObj == nil at start of ItemCarryGroundState, fallback to AirState")
				p.SetState(AirState)
				return
			}
			if p.HeldObj.Space.Label > 0 {
				// disbles collision
				p.HeldObj.UpdateLabel(-p.HeldObj.Space.Label)
			}
		},
		Loop: func(p *Player) {
			p.DoGroundCtrls()
			ItemCarryLoop(p)
			if p.Mods["jump"].JustActivated() {
				p.Jump()
				p.SetState(ItemCarryAirState)
			}
			p.DoGravity()
			if p.ActiColls.GroundHit == false {
				p.SetState(ItemCarryAirState)
			}
		},
	}.denil()
	ItemCarryAirState = PlayerState{
		Loop: func(p *Player) {
			if p.ActiColls.GroundHit {
				p.SetState(ItemCarryGroundState)
			}
			p.DoAirControls()
			p.DoGravity()
			ItemCarryLoop(p)
		},
	}.denil()
}
