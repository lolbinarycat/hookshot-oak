package player

import (
	"time"

	"github.com/lolbinarycat/hookshot-oak/direction"
	"github.com/lolbinarycat/hookshot-oak/labels"
	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/dlog"
)

//StateCommon is the function for commands that should be run in most
//states, like activating the hookshot if the hookshot button is pressed.
//It is not a state, and should not be used as one
func (p *Player) StateCommon() {
	p.IfHsPressedStartHs()
}

const FlySpeed = 4

var FlyState = PlayerState{
	Loop: func(p *Player) {
		dir := p.Ctrls.GetDir()
		p.Body.Delta.SetPos(dir.HCoeff()*FlySpeed,dir.VCoeff()*FlySpeed)
	},
}.denil()

var AirState PlayerState

func AirStateLoop(p *Player) {
	if p.PhysObject.ActiColls.GroundHit {
		p.SetState(GroundState)
		return
	} else {
		if p.Mods["walljump"].Active() && p.ActiColls.HLabel != labels.NoWallJump {
			if p.PhysObject.ActiColls.LeftWallHit {
				p.SetState(WallSlideLeftState)
			} else if p.PhysObject.ActiColls.RightWallHit {
				p.SetState(WallSlideRightState)
			}
		}

		p.DoGravity()
	}

	if p.Ctrls.GetDir().IsDown() && p.Mods["groundpound"].Active() {
		dlog.Info("groundpound started")
		p.SetState(GroundPoundStartState)
	}

	p.DoAirControls()
	p.StateCommon()
}

var RespawnFallState = PlayerState{
	Loop: func(p *Player) {
		if p.PhysObject.ActiColls.GroundHit {
			p.SetState(GroundState)
			return
		}
		p.DoGravity()
	},
}.denil()

var GroundState PlayerState

func GroundStateLoop(p *Player) {

	if p.PhysObject.ActiColls.GroundHit == true {
		if p.IsJumpInput() {
			p.Jump()
		}

	} else {
		p.SetState(CoyoteState)
	}

	if p.ActiColls.HLabel == labels.Block && p.ActiColls.VLabel != labels.Block {
		if p.Mods["blockpush"].Active() {
			if p.ActiColls.RightWallHit {
				p.SetState(BlockPushState(false))
			} else if p.ActiColls.LeftWallHit {
				p.SetState(BlockPushState(true))
			}
		}
	}

	if p.Mods["xdash"].JustActivated() {
		p.SetState(XDashState)
	}

	p.DoGroundCtrls()
	p.DoGravity()
	p.StateCommon()
}

func (p *Player) DoGroundCtrls() {
	p.Body.Delta.SetX(p.HeldDir.HCoeff()*p.Body.Speed.X())
}

const BlockPushSpeed float64 = 1

func BlockPushState(isLeft bool) PlayerState {
	return PlayerState{
		Start: func(p *Player) {
			if isLeft {
				p.GrabObjLeft(labels.Block)
			} else {
				p.GrabObjRight(labels.Block)
			}
		},
		Loop: func(p *Player) {
			if p.HeldObj == nil {
				p.SetState(GroundState)
			}
			if (isLeft && oak.IsDown(currentControls.Left) == false) ||
				(isLeft == false && oak.IsDown(curCtrls.Right) == false) {
				p.HeldObj.Delta.SetX(0)
				p.SetState(GroundState)
				return
			} else {
				var spd float64
				if isLeft {
					spd = -BlockPushSpeed
				} else {
					spd = BlockPushSpeed
				}

				p.Body.Delta.SetX(spd)
				p.HeldObj.Delta.SetX(spd)

				p.StateCommon()
			}
		},
	}.denil()
}

const BlockPullSpeed float64 = BlockPushSpeed

var BlockPullRightState = PlayerState{
	Loop: func(p *Player) {
		//if either button isn't pushed
		if (p.Mods["hs"].JustActivated() && oak.IsDown(currentControls.Right)) == false {
			p.SetState(GroundState)
			return
		}

		p.Body.Delta.SetX(BlockPullSpeed)
		p.HeldObj.Delta.SetX(BlockPullSpeed)
	},
}.denil()

//JumpHeightDecTime is how long JumpHeightDecState lasts
const JumpHeightDecTime time.Duration = time.Millisecond * 150

//the function JumpHeightDecState is the state that decides the height of the players jump.
//it does this by decreasing the gravity temporaraly when jump is held.
const MinHeightJumpInputTime = time.Millisecond * 85

var JumpHeightDecState = PlayerState{
	Loop: func(p *Player) {
		if p.TimeFromStateStart() > JumpHeightDecTime {
			p.SetState(AirState)
			return
		}
		if p.TimeFromStateStart() < MinHeightJumpInputTime &&
			!p.Mods["jump"].Active() {
			p.Body.Delta.SetY(-float64(JumpHeight) / 2)
			p.SetState(AirState)
			return
		}
		if p.Mods["jump"].Active() {
			//p.DoCustomGravity(Gravity/5)
		} else {
			p.DoGravity()
		}
		p.DoAirControls()
		p.StateCommon()
	},
}.denil()

//CoyoteTime is how long CoyoteState lasts
const CoyoteTime time.Duration = time.Millisecond * 7

//CoyoteState implements "coyote time" a window of time after
//running off an edge in which you can still jump
var CoyoteState = PlayerState{
	Loop: func(p *Player) {
		if p.StateStartTime.Add(CoyoteTime).Before(time.Now()) {
			p.SetState(AirState)
		}
		//inherit code from AirState
		//p.AirState()

		if p.ActiColls.RightWallHit && p.ActiColls.HLabel == labels.Block {
			p.SetState(BlockPushState(false))
		} else if p.ActiColls.LeftWallHit && p.ActiColls.HLabel == labels.Block {
			p.SetState(BlockPushState(true))
		}

		if p.IsJumpInput() {
			p.Jump()
		}

		p.StateCommon()
	},
}.denil()

var WallSlideLeftState = PlayerState{
	Loop: func(p *Player) {
		if p.Mods["climb"].Active() {
			p.SetState(ClimbLeftState)
			return
		}
		if p.IsJumpInput() {
			p.WallJump(direction.MaxRight(), true)
			return
		}
		if p.ActiColls.LeftWallHit == false {
			p.SetState(AirState)
			return
		}
		AirState.Loop(p)
	},
}.denil()

var WallSlideRightState = PlayerState{
	Loop: func(p *Player) {
		if p.Mods["climb"].Active() {
			p.SetState(ClimbRightState)
			return //return to stop airstate for overwriting our change
		}
		if p.IsJumpInput() {
			p.WallJump(direction.MaxLeft(), true)
			return
		}
		if p.ActiColls.RightWallHit == false {
			p.SetState(AirState)
			return
		}
		AirState.Loop(p)
	},
}.denil()

const WallJumpLaunchDuration time.Duration = time.Millisecond * 230

//func WallJumpLaunchState is entered after you walljump,
//temporaraly disabling your controls. This should prevent one sided
//walljumps
var WallJumpLaunchState = PlayerState{
	Loop: func(p *Player) {
		if p.TimeFromStateStart() >= WallJumpLaunchDuration {
			p.SetState(AirState)
			return
		}
		p.DoGravity()
		p.StateCommon()
	},
}.denil()

var ClimbRightState = PlayerState{
	Loop: func(p *Player) {
		if p.IsJumpInput() {
			p.WallJump(direction.MaxRight(), oak.IsDown(currentControls.Left))
			return
		}
		p.DoCliming()
		p.Body.Delta.SetX(1)
		//don't call StateCommon() here because it is called in DoCliming
	},
}.denil()

var ClimbLeftState = PlayerState{
	Loop: func(p *Player) {
		if p.IsJumpInput() {
			p.WallJump(direction.MaxRight(), oak.IsDown(currentControls.Right))
			return
		}
		p.Body.Delta.SetX(-1)
		p.DoCliming()
		//don't call StateCommon() because ... (see ClimbRightState)
	},
}.denil()

var ItemCarryGroundState PlayerState
var ItemCarryAirState PlayerState

func ItemCarryLoop(p *Player) {
	if p.Mods["hs"].JustActivated() {
		p.ThrowHeldItem(5*p.HeldDir.HCoeff(), -7)
		p.SetState(ItemThrowLag)
	} else {
		p.HeldObj.SetPos(p.Body.X(), p.Body.Y()-p.HeldObj.H)
		p.HeldObj.Delta.SetPos(p.Body.Delta.GetPos())
		p.HeldObj.ShiftPos(p.HeldObj.Delta.GetPos())
	}
}

func (p *Player) ThrowHeldItem(xSpeed, ySpeed float64) {
	p.HeldObj.Delta.SetPos(xSpeed, ySpeed)
	if p.HeldObj.Space.Label < 0 {
		p.HeldObj.Space.UpdateLabel(-p.HeldObj.Space.Label)
	}
	p.HeldObj = nil
}

var ItemThrowLag = PlayerState{
	MaxDuration: Frame * 3,
	NextState:   &AirState,
}.denil()
