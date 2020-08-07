package player

import (
	"time"
	"math"

	"github.com/lolbinarycat/hookshot-oak/direction"
	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/lolbinarycat/hookshot-oak/player/condition"
	"github.com/lolbinarycat/hookshot-oak/physobj"
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

func (p *Player) HandleJump() {
	if p.IsJumpInput() {
		if p.Mods["longjump"].Active() && p.HeldDir.IsDown() {
			p.Delta.SetY(-LongJumpH)
			p.Delta.SetX(p.Delta.X()/2+LongJumpW*p.HeldDir.HCoeff())
			p.SetState(AirState)
			dlog.Info("longjump")
			return
		}
		dlog.Info(p.HeldDir)
		p.Jump()
	}
}

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

	if p.HeldDir.IsDown() && !p.LastHeldDir.IsDown() &&
		p.Mods["groundpound"].Active() {

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

const LongJumpH = 4
const LongJumpW = 10

func GroundStateLoop(p *Player) {

	if p.PhysObject.ActiColls.GroundHit == true {
		p.HandleJump()
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

const (
	luigiMode_groundAccel = 1
	luigiMode_groundFriction = 1.17
)

func (p *Player) DoGroundCtrls() {
	if p.Mods["luigi"].Active() {
		p.Delta.SetX(p.Delta.X()/luigiMode_groundFriction)
		p.Delta.ShiftX(p.HeldDir.HCoeff()*luigiMode_groundAccel)
	} else {
		if math.Abs(p.Delta.X()) > math.Abs(p.Speed.X()) {
			p.Delta.SetX(p.Delta.X()*0.9)
		} else {
			p.Body.Delta.SetX(p.HeldDir.HCoeff()*p.Speed.X())
		}
		// TODO: make this feature into it's own module with
		// it's own states. ("crawl"?)
		if p.HeldDir.IsDown() {
			p.Delta.SetX(p.Delta.X()*0.6)
		}
	}
}

const BlockPushSpeed float64 = 1

// isHeldObjectNil checks if the player is not holding an object.
// it exists as a variable so that it's address can be taken,
// which needs to be done so it can be used as a map key
var isHeldObjectNil condition.Func = func (p interface{}) bool {
	return p.(*Player).HeldObj == nil 
}


func BlockPushState(isLeft bool) PlayerState {
	return PlayerState{
		//Map: map[condition.Condition]PlayerStateMapFunc{
		//	&isHeldObjectNil:constState(&AirState),
		//},
		Start: func(p *Player) {
			lastHit := p.ActiColls.LastHitH.E().(*physobj.Block)
			p.HeldObj = lastHit
			//runtime.Breakpoint()
			if p.HeldObj == nil {
				//lastHit := p.ActiColls.LastHitH.E()
				dlog.Error("Could not enter BlockPushState: HeldObj is nil")
				//p.HeldObj = lastHit.(*entities.Moving)
				//return
				p.SetState(GroundState)
				return
			}
		},
		Loop: func(p *Player) {

			if (isLeft && oak.IsDown(currentControls.Left) == false) ||
				(isLeft == false && oak.IsDown(curCtrls.Right) == false) {
				p.HeldObj.Body.Delta.SetX(0)
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
				p.HeldObj.Body.Delta.SetX(spd)

				p.StateCommon()
			}
		},
	}.denil()
}

const BlockPullSpeed float64 = BlockPushSpeed

// var BlockPullRightState = PlayerState{
// 	Loop: func(p *Player) {
// 		//if either button isn't pushed
// 		if (p.Mods["hs"].JustActivated() && oak.IsDown(currentControls.Right)) == false {
// 			p.SetState(GroundState)
// 			return
// 		}

// 		p.Body.Delta.SetX(BlockPullSpeed)
// 		p.HeldObj.Delta.SetX(BlockPullSpeed)
// 	},
// }.denil()

//JumpHeightDecTime is how long JumpHeightDecState lasts
const JumpHeightDecTime time.Duration = time.Millisecond * 150

//the function JumpHeightDecState is the state that decides the height of the players jump.
//it does this by decreasing the gravity temporaraly when jump is held.
const MinHeightJumpInputTime = time.Millisecond * 85

var JumpHeightDecState = PlayerState{
	Map: map[condition.Condition]PlayerStateMapFunc{
		&condition.FramesElapsed{N:14}:func(*Player) *PlayerState {
			return &AirState
		},
		&condition.FramesElapsed{N:4}:func(p *Player) *PlayerState {
			if !p.Mods["jump"].Active() {
				p.Body.Delta.SetY(-float64(JumpHeight) / 2)
				return &AirState
			}
			// continue evaluating map, or if there are no more elements, eval Loop
			return nil
		},
	},
	Loop: func(p *Player) {
		if p.Mods["jump"].Active() {
			//p.DoCustomGravity(Gravity/5)
		} else {
			p.DoGravity()
		}
		p.DoAirControls()
		p.StateCommon()
	},
}.denil()

const CoyoteFrames = 7

//CoyoteState implements "coyote time" a window of time after
//running off an edge in which you can still jump
var CoyoteState = PlayerState{
	Map: map[condition.Condition]PlayerStateMapFunc{
		&condition.FramesElapsed{N:CoyoteFrames}:constState(&AirState),
	},
	Loop: func(p *Player) {
		if p.PhysObject.ActiColls.GroundHit == true {
			p.SetState(GroundState)
		}

		if p.ActiColls.RightWallHit && p.ActiColls.HLabel == labels.Block {
			p.SetState(BlockPushState(false))
		} else if p.ActiColls.LeftWallHit && p.ActiColls.HLabel == labels.Block {
			p.SetState(BlockPushState(true))
		}

		p.HandleJump()

		p.DoGravity()
		p.DoAirControls()
		p.StateCommon()
	},
}.denil()

func constState(state *PlayerState) PlayerStateMapFunc {
	return func (_ *Player) *PlayerState {
		return state
	}
}

var WallSlideLeftState = PlayerState{
	LLoop: func(p *Player) *PlayerState {
		if p.Mods["climb"].Active() {
			return &ClimbLeftState
		}
		if p.ActiColls.LeftWallHit == false {
			return &AirState
		}
		return nil
	},
	Loop: func(p *Player) {
		if p.IsJumpInput() {
			p.WallJump(direction.MaxRight(), true)
			return
		}

		AirState.Loop(p)
	},
}.denil()

var WallSlideRightState = PlayerState{
	LLoop: func(p *Player) *PlayerState {
		if p.Mods["climb"].Active() {
			return &ClimbRightState
		}
		if p.ActiColls.RightWallHit == false {
			return &AirState
		}
		return nil
	},
	Loop: func(p *Player) {
		if p.IsJumpInput() {
			p.WallJump(direction.MaxLeft(), true)
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
		p.HeldObj.Body.SetPos(p.Body.X(), p.Body.Y()-p.HeldObj.Body.H)
		p.HeldObj.Body.Delta.SetPos(p.Body.Delta.GetPos())
		p.HeldObj.Body.ShiftPos(p.HeldObj.Delta.GetPos())
	}
}

func (p *Player) ThrowHeldItem(xSpeed, ySpeed float64) {
	p.HeldObj.Body.Delta.SetPos(xSpeed, ySpeed)
	p.HeldObj.Held = false
	if p.HeldObj.Body.Space.Label < 0 {
		p.HeldObj.Body.Space.UpdateLabel(-p.HeldObj.Space.Label)
	}
	p.HeldObj = nil
}

var ItemThrowLag = PlayerState{
	Map: map[condition.Condition]PlayerStateMapFunc{
		&condition.FramesElapsed{N:3}:constState(&AirState),
	},
	NextState:   &AirState,
}.denil()
