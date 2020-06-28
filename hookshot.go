package main

import (
	"time"

	"github.com/lolbinarycat/hookshot-oak/labels"
	oak "github.com/oakmound/oak/v2"
)

// This file contains functions that handle hookshot behavior

func (p *Player) HsItemGrabLoop(dir Direction) {
	if (dir == Right && (p.Hs.X <= 0 || p.ActiColls.RightWallHit)) ||
		(dir == Left && (p.Hs.X >= 0 || p.ActiColls.LeftWallHit)) {
		p.EndHs()
		return
	}


	var coeff float64
	if dir == Right {
		coeff = -1
	} else if dir == Left {
		coeff = 1
	} else {
		panic("unknown direction given to HsItemGrabLoop")
	}

	p.Hs.Body.Delta.SetX(p.Hs.Body.Speed.X() * coeff)
	p.HeldObj.Delta.SetX(p.Hs.Body.Delta.X())
}

func (p *Player) HsExtendLoop(dir Direction) {
	p.Hs.Active = true
	if p.TimeFromStateStart() > HsExtendTime {
		if dir == Right {
			p.SetState(HsRetractRightState)
		} else if dir == Left {
			p.SetState(HsRetractLeftState)
		} else {
			panic("unknown direction given to HsExtendLoop")
		}
	} else if dir == Right && p.Hs.ActiColls.RightWallHit {
		if p.Hs.ActiColls.HLabel == labels.Block && p.Mods.HsItemGrab.Equipped {
			p.SetState(HsItemGrabRightState)
			return
		} else {
			p.SetState(HsPullRightState)
			return
		}
	} else if dir == Left && p.Hs.ActiColls.LeftWallHit {
		if p.Hs.ActiColls.HLabel == labels.Block && p.Mods.HsItemGrab.Equipped {
			p.SetState(HsItemGrabLeftState)
			return
		} else {
			p.SetState(HsPullRightState)
			return
		}
	} else if p.TimeFromStateStart() > HsInputTime && isHsInput() {
			if dir == Right {
				p.SetState(HsRetractRightState)
			} else if dir == Left {
				p.SetState(HsRetractLeftState)
			} else {
				panic("invalid direction to HsExtendLoop")
			}
			return
	} else {
		p.Body.Delta.SetPos(0, 0)
		if dir == Right {
			p.Hs.Body.Delta.SetX(p.Hs.Body.Speed.X())
		} else if dir == Left {
			p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
		}
	}
}


//STATES:


const HsStartTime time.Duration = time.Millisecond * 60

var HsStartState = PlayerState{
	Loop: func(p *Player) {
		if player.Mods.Hookshot.Equipped == false {
			p.SetState(AirState)
			return
		}
		if p.TimeFromStateStart() > HsStartTime {
			if oak.IsDown(currentControls.Right) {
				p.SetState(HsExtendRightState)
			} else if oak.IsDown(currentControls.Left) {
				p.SetState(HsExtendLeftState)
			} else {
				p.SetState(AirState)
			}
		}
	}}.denil()

var HsExtendRightState = PlayerState{
	Loop: func(p *Player) {
		p.HsExtendLoop(Right)
	},
}.denil()

var HsExtendLeftState = PlayerState{
	Loop: func(p *Player) {
		p.HsExtendLoop(Left)
	},
}.denil()

var HsRetractRightState = PlayerState{
	Loop: func(p *Player) {
		if p.Hs.X <= 0 {
			p.EndHs()
			return
		}
		p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
		//p.Hs.X -= p.Hs.Body.Speed.X()
	},
}.denil()

var HsRetractLeftState = PlayerState{
	Loop: func(p *Player) {
		if p.Hs.X >= 0 {
			p.EndHs()
			return
		}
		p.Hs.Body.Delta.SetX(p.Hs.Body.Speed.X())
		//p.Hs.X -= p.Hs.Body.Speed.X()
	},
}.denil()

//HsPullRightState is the state for when the hookshot is
//pulling the player after having hit an object
var HsPullRightState = PlayerState{
	Loop: func(p *Player) {
		if p.ActiColls.RightWallHit {
			p.EndHs()
			return
		}
		//p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
		p.Body.Delta.SetX(p.Hs.Body.Speed.X())
		//p.PullPlayer()
	},
}

var HsPullLeftState = PlayerState{
	Loop: func(p *Player) {
		if p.ActiColls.LeftWallHit {
			p.EndHs()
			return
		}
		//p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
		p.Body.Delta.SetX(-p.Hs.Body.Speed.X())
		//p.PullPlayer()
	},
}.denil()

var HsItemGrabRightState = PlayerState{
	Start: func(p *Player) {
		p.HeldObj = p.Hs.GetLastHitObj(true)
	},
	Loop: func(p *Player) {
		p.HsItemGrabLoop(Right)
	},
	End: func(p *Player) {
		p.HeldObj = nil
		p.HeldObj.Delta.SetPos(0,0)
	},
}.denil()

var HsItemGrabLeftState = PlayerState{
	Start: func(p *Player) {
		p.HeldObj = p.Hs.GetLastHitObj(true)
	},
	Loop: func(p *Player) {
		p.HsItemGrabLoop(Left)
	},
	End: func(p *Player) {
		p.HeldObj.Delta.SetPos(0,0)
		p.HeldObj = nil
	},
}.denil()
