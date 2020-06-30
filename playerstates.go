package main

import (
	"time"

	"github.com/lolbinarycat/hookshot-oak/direction"
	"github.com/lolbinarycat/hookshot-oak/labels"
	oak "github.com/oakmound/oak/v2"
)

//StateCommon is the function for commands that should be run in most
//states, like activating the hookshot if the hookshot button is pressed.
//It is not a state, and should not be used as one
func (p *Player) StateCommon() {
	p.ifHsPressedStartHs()
}

const FlySpeed = 4

var FlyState = PlayerState{
	Loop: func(p *Player) {
		if oak.IsDown(currentControls.Up) {
			p.Body.Delta.SetY(-FlySpeed)
		} else if oak.IsDown(curCtrls.Down) {
			p.Body.Delta.SetY(FlySpeed)
		} else {
			p.Body.Delta.SetY(0)
		}

		if oak.IsDown(curCtrls.Left) {
			p.Body.Delta.SetX(-FlySpeed)
		} else if oak.IsDown(curCtrls.Right) {
			p.Body.Delta.SetX(FlySpeed)
		} else {
			p.Body.Delta.SetX(0)
		}
	},
}.denil()


var AirState PlayerState

func AirStateLoop(p *Player) {

	if player.PhysObject.ActiColls.GroundHit {
		p.SetState(GroundState)
		return
	} else {
		if p.Mods.WallJump.Equipped && p.ActiColls.HLabel != labels.NoWallJump {
			if p.PhysObject.ActiColls.LeftWallHit {
				p.SetState(WallSlideLeftState)
			} else if p.PhysObject.ActiColls.RightWallHit {
				p.SetState(WallSlideRightState)
			}
		}

		p.DoGravity()
	}

	if oak.IsDown(curCtrls.Down) && p.Mods.GroundPound.Equipped {
		p.SetState(GroundPoundStartState)
	}

	p.DoAirControls()
	p.StateCommon()
}

var RespawnFallState = PlayerState{
	Loop: func(p *Player) {
		if player.PhysObject.ActiColls.GroundHit {
			p.SetState(GroundState)
			return
		}
		p.DoGravity()
	},
}.denil()

const GroundPoundStartTime = time.Second/5
var GroundPoundStartState = PlayerState{
	Start: func(p *Player) {
p.Body.Delta.SetPos(0,0)
	},
	Loop: func(p *Player) {
		if p.TimeFromStateStart() > GroundPoundStartTime {
			p.SetState(GroundPoundState)
		}
	},
}.denil()

const GroundPoundSpeed = 8
var GroundPoundState = PlayerState{
	Loop: func(p *Player) {
		if p.ActiColls.GroundHit {
			p.SetState(GroundState)
		} else if oak.IsDown(curCtrls.Up){
			p.Body.Delta.SetY(0)
			p.SetState(AirState)
		} else {
			p.Body.Delta.SetY(GroundPoundSpeed)
		}
	},
}.denil()

var GroundState PlayerState

func GroundStateLoop(p *Player) {

	if player.PhysObject.ActiColls.GroundHit == true {
		if isJumpInput() {
			p.Jump()
		}

	} else {
		p.SetState(CoyoteState)
	}

	if p.ActiColls.HLabel == labels.Block {
		if p.Mods.BlockPush.Equipped {
			if p.ActiColls.RightWallHit {
				p.SetState(BlockPushState(false))
			} else if p.ActiColls.LeftWallHit {
				p.SetState(BlockPushState(true))
			}
		}
	}

	if oak.IsDown(currentControls.Left) {
		player.Body.Delta.SetX(-player.Body.Speed.X())
	} else if oak.IsDown(currentControls.Right) {
		player.Body.Delta.SetX(player.Body.Speed.X())
	} else {
		player.Body.Delta.SetX(0)
		//player.Body.Delta.X()/2)
	}
	p.DoGravity()
	p.StateCommon()

}

const BlockPushSpeed float64 = 1

func BlockPushState(isLeft bool) PlayerState {
	return PlayerState{
	Start: func(p *Player) { p.GrabObjLeft(labels.Block) },
	Loop: func(p *Player) {
		if (isLeft && oak.IsDown(currentControls.Left) == false) ||
			(isLeft == false && oak.IsDown(curCtrls.Right) == false ) {
			p.HeldObj.Delta.SetX(0)
			p.SetState(GroundState)
		} else {
			var spd float64
			if isLeft {
				spd = -BlockPushSpeed
			} else {
				spd = BlockPushSpeed
			}

			p.Body.Delta.SetX(spd)
			p.HeldObj.Delta.SetX(spd)
		}
	},
	}.denil()
}

const BlockPullSpeed float64 = BlockPushSpeed

var BlockPullRightState = PlayerState{
	Loop: func(p *Player) {
		//if either button isn't pushed
		if (oak.IsDown(currentControls.Climb) && oak.IsDown(currentControls.Right)) == false {
			p.SetState(GroundState)
			return
		}

		p.Body.Delta.SetX(BlockPullSpeed)
		p.HeldObj.Delta.SetX(BlockPullSpeed)
	},
}.denil()

//the function JumpHeightDecState is the state that decides the height of the players jump.
//it does this by decreasing the gravity temporaraly when jump is held.
var JumpHeightDecState = PlayerState{
	Loop: func(p *Player) {
		if p.TimeFromStateStart() > JumpHeightDecTime {
			p.SetState(AirState)
			return
		}
		if oak.IsDown(currentControls.Jump) {
			//p.DoCustomGravity(Gravity/5)
		} else {
			p.DoGravity()
		}
		p.DoAirControls()
		p.StateCommon()
	},
}.denil()

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
			p.SetState(BlockPushRightState)
		} else if p.ActiColls.LeftWallHit && p.ActiColls.HLabel == labels.Block {
			p.SetState(BlockPushState(true))
		}

		if isJumpInput() {
			p.Jump()
		}

		p.StateCommon()
	},
}.denil()

var WallSlideLeftState = PlayerState{
	Loop: func(p *Player) {
		if oak.IsDown(currentControls.Climb) && p.Mods.Climb.Equipped {
			p.SetState(ClimbLeftState)
			return
		}
		if isJumpInput() {
			p.WallJump(direction.MaxLeft(), true)
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
		if oak.IsDown(currentControls.Climb) && p.Mods.Climb.Equipped {
			p.SetState(ClimbRightState)
			return //return to stop airstate for overwriting our change
		}
		if isJumpInput() {
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
		if isJumpInput() {
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
		if isJumpInput() {
			p.WallJump(direction.MaxRight(), oak.IsDown(currentControls.Right))
			return
		}
		p.Body.Delta.SetX(-1)
		p.DoCliming()
		//don't call StateCommon() because ... (see ClimbRightState)
	},
}.denil()

