package player

import (
	"time"
)

const GroundPoundStartTime = time.Second / 5

var GroundPoundStartState = PlayerState{
	Start: func(p *Player) {
		p.Body.Delta.SetPos(0, 0)
	},
	Loop: func(p *Player) {
		if p.TimeFromStateStart() > GroundPoundStartTime {
			p.SetState(GroundPoundState)
		}
		p.StateCommon()
	},
}.denil()

const GroundPoundSpeed = 8

var GroundPoundState = PlayerState{
	Loop: func(p *Player) {
		if p.ActiColls.GroundHit {
			p.SetState(GroundPoundEndState)
		} else if p.Ctrls.GetDir().IsUp() {
			p.Body.Delta.SetY(0)
			p.SetState(AirState)
		} else {
			p.Body.Delta.SetY(GroundPoundSpeed)
			p.StateCommon()
		}
	},
}.denil()

const GroundPoundEndTime = time.Millisecond * 80

var GroundPoundEndState = PlayerState{
	Loop: func(p *Player) {
		if p.TimeFromStateStart() > GroundPoundEndTime {
			p.SetState(GroundState)
		} else if p.IsJumpInput() && p.Mods["groundpoundjump"].Active() {
			p.SetState(GroundPoundJumpState)
		} else {
			p.Body.Delta.SetX(0)
		}
	},
	MaxDuration: GroundPoundEndTime,
	NextState:   &GroundState,
}.denil()

const GroundPoundJumpGravity float64 = Gravity / 2
const GroundPoundJumpTime = time.Millisecond * 60
const GroundPoundJumpForce = 9

var GroundPoundJumpState = PlayerState{
	Start: func(p *Player) {
		if p.Mods["groundpoundjump"].Active() == false {
			p.SetState(AirState)
		} else {
			p.Body.Delta.SetY(-GroundPoundJumpForce)
		}
	},
	Loop: func(p *Player) {
		if p.Body.Delta.Y() >= 0 || p.TimeFromStateStart() > GroundPoundJumpTime {
			p.SetState(AirState)
		} else {
			p.DoCustomGravity(GroundPoundJumpGravity)
			p.DoAirControls()
		}
	},
}.denil()
