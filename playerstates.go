package main

import (
	"time"

	"github.com/lolbinarycat/hookshot-oak/labels"
	oak "github.com/oakmound/oak/v2"
)

//StateCommon is the function for commands that should be run in most
//states, like activating the hookshot if the hookshot button is pressed.
//It is not a state, and should not be used as one 
func (p *Player) StateCommon() {
	p.ifHsPressedStartHs()
}


var AirState = PlayerState{
	Loop:func(p *Player) { //start in air state

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

	p.DoAirControls()
	p.StateCommon()
	},
	Start:func(p *Player) {},
	End:func(p *Player) {},
}

var RespawnFallState = PlayerState{
	Loop:func (p *Player)  {
		if player.PhysObject.ActiColls.GroundHit {
			p.SetState(GroundState)
			return
		}
		p.DoGravity()
	},
}

var  GroundState = PlayerState{
	Loop:func(p *Player) {

	if player.PhysObject.ActiColls.GroundHit == true {
		if isJumpInput() {
			p.Jump()
		}

	} else {
		p.SetState(CoyoteState)
	}

	if  p.ActiColls.HLabel == labels.Block {
		if p.Mods.BlockPush.Equipped {
			if p.ActiColls.RightWallHit  {
				p.GrabObj(labels.Block) //temporay, until state overhaul
				p.SetState(BlockPushRightState)
			} else if p.ActiColls.LeftWallHit {
				p.SetState(BlockPushLeftState)
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
},
}

const BlockPushSpeed float64 = 1
var BlockPushRightState = PlayerState{
	Loop:func (p *Player)  {
	if oak.IsDown(currentControls.Right) == false {
		block.Body.Delta.SetX(0)
		p.SetState(p.GroundState)
		return
	}
	p.Body.Delta.SetX(BlockPushSpeed)
	
	//if id > 0 {
		//id.E().(*entities.Moving).Delta.SetX(BlockPushSpeed)
		p.HeldObj.Delta.SetX(BlockPushSpeed)
	//}
	

	//p.HeldObj.Delta.ShiftX(BlockPushSpeed)
	//hitBlock := p.Body.HitLabel(labels.Block)
	//block.Body.Delta.SetX(BlockPushSpeed)

	
	/*if blk, ok := event.GetEntity(2).(*entities.Moving); ok && blk == block.Body {
		blk.Delta.SetX(BlockPushSpeed)

	} else {
		dlog.Error("type assertion failed")
	}*/
	
	//if hitBlock != nil {
	//}
},
}

var BlockPushLeftState = PlayerState{
	Loop: func (p *Player)  {
	if oak.IsDown(currentControls.Left) == false {
		block.Body.Delta.SetX(0)
		p.SetState(p.GroundState)
		return
	}
	p.Body.Delta.SetX(-BlockPushSpeed)
	//hitBlock := p.Body.HitLabel(labels.Block)
	block.Body.Delta.SetX(-BlockPushSpeed)

},
}

const BlockPullSpeed float64 = BlockPushSpeed
var BlockPullRightState = PlayerState{
Loop:func (p *Player)  {
	//if either button isn't pushed
	if (oak.IsDown(currentControls.Climb) && oak.IsDown(currentControls.Right)) == false {
		p.SetState(p.GroundState)
		return
	}

	p.Body.Delta.SetX(BlockPullSpeed)
	block.Body.Delta.SetX(BlockPullSpeed)
}
}

//the function JumpHeightDecState is the state that decides the height of the players jump.
//it does this by decreasing the gravity temporaraly when jump is held.
var JumpHeightDecState = PlayerState{
Loop:func (p *Player)  {
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
},
}

//CoyoteState implements "coyote time" a window of time after
//running off an edge in which you can still jump
var  CoyoteState = PlayerState{
Loop:func (p *Player) {
	if p.StateStartTime.Add(CoyoteTime).Before(time.Now()) {
		p.SetState(p.AirState)
	}
	//inherit code from AirState
	//p.AirState()

	if p.ActiColls.RightWallHit && p.ActiColls.HLabel == labels.Block {
		p.SetState(p.BlockPushRightState)
	} else if p.ActiColls.LeftWallHit && p.ActiColls.HLabel == labels.Block {
		p.SetState(p.BlockPushLeftState)
	}

	if isJumpInput() {
		p.Jump()
	}

	p.StateCommon()
},
}

var WallSlideLeftState = PlayerState{
	Loop:func (p *Player)  {
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
	},
}

var  WallSlideRightState = PlayerState{
	Loop:func (p *Player) {
	if oak.IsDown(currentControls.Climb) && p.Mods.Climb.Equipped {
		p.SetState(ClimbRightState)
		return //return to stop airstate for overwriting our change
	}
	if isJumpInput() {
		p.WallJump(Left, true)
		return
	}
	if p.ActiColls.RightWallHit == false {
		p.SetState(AirState)
		return
	}
	p.AirState()
},
}

//func WallJumpLaunchState is entered after you walljump,
//temporaraly disabling your controls. This should prevent one sided
//walljumps
var WallJumpLaunchState = PlayerState{
	Loop:func (p *Player)  {
	if p.TimeFromStateStart() >= WallJumpLaunchDuration {
		p.SetState(p.AirState)
		return
	}
	p.DoGravity()
	p.StateCommon()
},
}

var  ClimbRightState = PlayerState{
Loop:func (p *Player) {
	if isJumpInput() {
		p.WallJump(Left, oak.IsDown(currentControls.Left))
		return
	}
	p.DoCliming()
	p.Body.Delta.SetX(1)
	//don't call StateCommon() here because it is called in DoCliming
},
}

var ClimbLeftState = PlayerState{
func (p *Player) {
	if isJumpInput() {
		p.WallJump(Right, oak.IsDown(currentControls.Right))
		return
	}
	p.Body.Delta.SetX(-1)
	p.DoCliming()
	//don't call StateCommon() because ... (see ClimbRightState)
}
}
const HsStartTime time.Duration = time.Millisecond * 60
var HsStartState = PlayerState{
	Loop:func (p *Player)  {
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
}}

var HsExtendRightState = PlayerState{
	Loop:func (p *Player)  {
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
	},
}

var  HsExtendLeftState = PlayerState{
	Loop:func (p *Player) {
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
	},
}

var HsRetractRightState = PlayerState{
	Loop:func (p *Player)  {
	if p.Hs.X <= 0 {
		p.EndHs()
		return
	}
	p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
	//p.Hs.X -= p.Hs.Body.Speed.X()
}
}

var HsRetractLeftState = PlayerState{
	Loop:func (p *Player)  {
	if p.Hs.X >= 0 {
		p.EndHs()
		return
	}
	p.Hs.Body.Delta.SetX(p.Hs.Body.Speed.X())
	//p.Hs.X -= p.Hs.Body.Speed.X()
},
}

//HsPullRightState is the state for when the hookshot is
//pulling the player after having hit an object
var HsPullRightState = PlayerState{
Loop:func (p *Player) {
	if p.ActiColls.RightWallHit {
		p.EndHs()
		return
	}
	//p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
	p.Body.Delta.SetX(p.Hs.Body.Speed.X())
	//p.PullPlayer()
},
}

var  HsPullLeftState = PlayerState{
	Loop:func (p *Player) {
	if p.ActiColls.LeftWallHit {
		p.EndHs()
		return
	}
	//p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
	p.Body.Delta.SetX(-p.Hs.Body.Speed.X())
	//p.PullPlayer()
	},
}

