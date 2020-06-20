package main

import (
	"time"

	oak "github.com/oakmound/oak/v2"
)


func (p *Player) AirState() { //start in air state

	if player.PhysObject.ActiColls.GroundHit {
		p.SetState(p.GroundState)
		return
	} else {
		if p.Mods.WallJump.Equipped && p.ActiColls.HLabel != NoWallJump{
			if p.PhysObject.ActiColls.LeftWallHit {
				p.SetState(p.WallSlideLeftState)
			} else if p.PhysObject.ActiColls.RightWallHit {
				p.SetState(p.WallSlideRightState)
			}
		}

		p.DoGravity()
	}

	p.DoAirControls()
	p.ifHsPressedStartHs()
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
}

func (p *Player) ClimbRightState() {
	if isJumpInput() {
		p.WallJump(Left, oak.IsDown(currentControls.Left))
		return
	}
	p.DoCliming()
	//if p.Body.Space.Above(p.ActiColls.LastHitH)
	//p.Body.Delta.SetX(1)
}

func (p *Player) ClimbLeftState() {
	if isJumpInput() {
		p.WallJump(Right, oak.IsDown(currentControls.Right))
		return
	}
	p.Body.Delta.SetX(-1)
	p.DoCliming()
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
		} else {
			p.SetState(p.AirState)
		}
	}
}

func (p *Player) HsExtendRightState() {
	if p.TimeFromStateStart() > HsExtendTime {
		p.SetState(p.HsRetractRightState)
	} else if p.Hs.ActiColls.RightWallHit {
		p.SetState(p.HsRetractRightState)
	} else {
		p.Hs.Active = true
		p.Body.Delta.SetPos(0,0)
		p.Hs.Body.Delta.SetX(p.Hs.Body.Speed.X())
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
