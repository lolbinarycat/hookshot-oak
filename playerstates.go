package main

import (
	"time"

	oak "github.com/oakmound/oak/v2"
	"github.com/lolbinarycat/hookshot-oak/labels"
)

//StateCommon is the function for commands that should be run in most
//states, like activating the hookshot if the hookshot button is pressed.
//It is not a state, and should not be used as one 
func (p *Player) StateCommon() {
	p.ifHsPressedStartHs()
}

func (p *Player) AirState() { //start in air state

	if player.PhysObject.ActiColls.GroundHit {
		p.SetState(p.GroundState)
		return
	} else {
		if p.Mods.WallJump.Equipped && p.ActiColls.HLabel != labels.NoWallJump {
			if p.PhysObject.ActiColls.LeftWallHit {
				p.SetState(p.WallSlideLeftState)
			} else if p.PhysObject.ActiColls.RightWallHit {
				p.SetState(p.WallSlideRightState)
			}
		}

		p.DoGravity()
	}

	p.DoAirControls()
	p.StateCommon()
}

func (p *Player) RespawnFallState() {
	if player.PhysObject.ActiColls.GroundHit {
		p.SetState(p.GroundState)
		return
	}
	p.DoGravity()
}

func (p *Player) GroundState() {

	if player.PhysObject.ActiColls.GroundHit == true {
		if isJumpInput() {
			p.Jump()
		}

	} else {
		p.SetState(player.CoyoteState)
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

//the function JumpHeightDecState is the state that decides the height of the players jump.
//it does this by decreasing the gravity temporaraly when jump is held.
func (p *Player) JumpHeightDecState() {
	if p.TimeFromStateStart() > JumpHeightDecTime {
		p.SetState(p.AirState)
		return
	}
	if oak.IsDown(currentControls.Jump) {
		//p.DoCustomGravity(Gravity/5)
	} else {
		p.DoGravity()
	}
	p.DoAirControls()
	p.StateCommon()
}

//CoyoteState implements "coyote time" a window of time after
//running off an edge in which you can still jump
func (p *Player) CoyoteState() {
	if p.StateStartTime.Add(CoyoteTime).Before(time.Now()) {
		p.SetState(p.AirState)
	}
	//inherit code from AirState
	//p.AirState()
	if isJumpInput() {
		p.Jump()
	}

	p.StateCommon()
}

func (p *Player) WallSlideLeftState() {
	if oak.IsDown(currentControls.Climb) && p.Mods.Climb.Equipped {
		p.SetState(p.ClimbLeftState)
		return
	}
	if isJumpInput() {
		p.WallJump(Right, true)
		return
	}
	if p.ActiColls.LeftWallHit == false {
		p.SetState(p.AirState)
		return
	}
	p.AirState()
}

func (p *Player) WallSlideRightState() {
	if oak.IsDown(currentControls.Climb) && p.Mods.Climb.Equipped {
		p.SetState(p.ClimbRightState)
		return //return to stop airstate for overwriting our change
	}
	if isJumpInput() {
		p.WallJump(Left, true)
		return
	}
	if p.ActiColls.RightWallHit == false {
		p.SetState(p.AirState)
		return
	}
	p.AirState()
}

//func WallJumpLaunchState is entered after you walljump,
//temporaraly disabling your controls. This should prevent one sided
//walljumps
func (p *Player) WallJumpLaunchState() {
	if p.TimeFromStateStart() >= WallJumpLaunchDuration {
		p.SetState(p.AirState)
		return
	}
	p.DoGravity()
	p.StateCommon()
}

func (p *Player) ClimbRightState() {
	if isJumpInput() {
		p.WallJump(Left, oak.IsDown(currentControls.Left))
		return
	}
	p.DoCliming()
	p.Body.Delta.SetX(1)
	//don't call StateCommon() here because it is called in DoCliming
}

func (p *Player) ClimbLeftState() {
	if isJumpInput() {
		p.WallJump(Right, oak.IsDown(currentControls.Right))
		return
	}
	p.Body.Delta.SetX(-1)
	p.DoCliming()
	//don't call StateCommon() because ... (see ClimbRightState)
}

const HsStartTime time.Duration = time.Millisecond * 60
func (p *Player) HsStartState() {
	if player.Mods.Hookshot.Equipped == false {
		p.SetState(p.AirState)
		return
	}
	if p.TimeFromStateStart() > HsStartTime {
		if oak.IsDown(currentControls.Right) {
			p.SetState(p.HsExtendRightState)
		} else if oak.IsDown(currentControls.Left) {
			p.SetState(p.HsExtendLeftState)
		} else {
			p.SetState(p.AirState)
		}
	}
}

func (p *Player) HsExtendRightState() {
	p.Hs.Active = true
	if p.TimeFromStateStart() > HsExtendTime {
		p.SetState(p.HsRetractRightState)

	} else if p.Hs.ActiColls.RightWallHit {
		p.SetState(p.HsPullRightState)

	} else {
		if p.TimeFromStateStart() > HsInputTime && isHsInput() {
			p.SetState(p.HsRetractRightState)
			return
		}
		p.Body.Delta.SetPos(0,0)
		p.Hs.Body.Delta.SetX(p.Hs.Body.Speed.X())
	}
}

func (p *Player) HsExtendLeftState() {
	p.Hs.Active = true
	if p.TimeFromStateStart() > HsExtendTime {
		p.SetState(p.HsRetractLeftState)

	} else if p.Hs.ActiColls.LeftWallHit {
		p.SetState(p.HsPullLeftState)

	} else {
		if p.TimeFromStateStart() > HsInputTime && isHsInput() {
			p.SetState(p.HsRetractLeftState)
			return
		}
		p.Body.Delta.SetPos(0,0)
		p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
	}
}

func (p *Player) HsRetractRightState() {
	if p.Hs.X <= 0 {
		p.EndHs()
		return
	}
	p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
	//p.Hs.X -= p.Hs.Body.Speed.X()
}

func (p *Player) HsRetractLeftState() {
	if p.Hs.X >= 0 {
		p.EndHs()
		return
	}
	p.Hs.Body.Delta.SetX(p.Hs.Body.Speed.X())
	//p.Hs.X -= p.Hs.Body.Speed.X()
}

//HsPullRightState is the state for when the hookshot is
//pulling the player after having hit an object
func (p *Player) HsPullRightState() {
	if p.ActiColls.RightWallHit {
		p.EndHs()
		return
	}
	//p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
	p.Body.Delta.SetX(p.Hs.Body.Speed.X())
	//p.PullPlayer()
}

func (p *Player) HsPullLeftState() {
	if p.ActiColls.LeftWallHit {
		p.EndHs()
		return
	}
	//p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
	p.Body.Delta.SetX(-p.Hs.Body.Speed.X())
	//p.PullPlayer()
}

