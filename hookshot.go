package main

import (
	"time"

	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/lolbinarycat/hookshot-oak/direction"

	oak "github.com/oakmound/oak/v2"
)

// This file contains functions that handle hookshot behavior

func (p *Player) HsItemGrabLoop(dir direction.Dir) {
	if (dir.IsJustRight() && (p.Hs.X <= 0 || p.ActiColls.RightWallHit)) ||
		(dir.IsJustLeft() && (p.Hs.X >= 0 || p.ActiColls.LeftWallHit)) {
		p.EndHs()
		return
	}


	var coeff float64
	if dir.IsJustRight() {
		coeff = -1
	} else if dir.IsJustLeft(){
		coeff = 1
	} else {
		panic("unknown direction given to HsItemGrabLoop")
	}

	p.Hs.Body.Delta.SetX(p.Hs.Body.Speed.X() * coeff)
	p.HeldObj.Delta.SetX(p.Hs.Body.Delta.X())
}

func (p *Player) HsExtendLoop(dir direction.Dir) {
	p.Hs.Active = true
	if p.TimeFromStateStart() > HsExtendTime {
		if dir.IsJustRight() {
			p.SetState(HsRetractRightState)
		} else if dir.IsJustLeft() {
			p.SetState(HsRetractLeftState)
		} else {
			panic("unknown direction given to HsExtendLoop")
		}
	} else if dir.IsJustRight() && p.Hs.ActiColls.RightWallHit {
		if p.Hs.ActiColls.HLabel == labels.Block && p.Mods.HsItemGrab.Equipped {
			p.SetState(HsItemGrabRightState)
			return
		} else {
			p.SetState(HsPullState(dir))
			return
		}
	} else if dir.IsJustLeft() && p.Hs.ActiColls.LeftWallHit {
		if p.Hs.ActiColls.HLabel == labels.Block && p.Mods.HsItemGrab.Equipped {
			p.SetState(HsItemGrabLeftState)
			return
		} else {
			p.SetState(HsPullState(dir))
			return
		}
	} else if p.TimeFromStateStart() > HsInputTime && isHsInput() {
			if dir.IsJustRight() {
				p.SetState(HsRetractRightState)
			} else if dir.IsJustLeft() {
				p.SetState(HsRetractLeftState)
			} else {
				panic("invalid direction to HsExtendLoop")
			}
			return
	} else {
		p.Body.Delta.SetPos(0, 0)
		if dir.IsJustRight() {
			p.Hs.Body.Delta.SetX(p.Hs.Body.Speed.X())
		} else if dir.IsJustLeft() {
			p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
		}
	}
}

func (p *Player) HsPullLoop(dir direction.Dir) {
	panic("no longer used")
}

func HsPullState(dir direction.Dir) PlayerState {
	return PlayerState{
		Loop: func(p *Player) {
			if dir.IsRight() {
				if p.ActiColls.RightWallHit {
					p.EndHs()
					return
				}
			} else if dir.IsLeft() {
				if p.ActiColls.LeftWallHit {
					p.EndHs()
					return
				}
			}
			coeffX := direction.ToCoeff(dir.H)
			coeffY := direction.ToCoeff(dir.V)
			p.Body.Delta.SetPos(coeffX * p.Hs.Body.Speed.X(),
				coeffY * p.Hs.Body.Speed.Y())
		},
	}.denil()
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
			//} else if oak.IsDown(currentControls.Up) {
			//p.SetState(HsExtendUpState)
			} else {
				p.SetState(AirState)
			}
		}
	}}.denil()

var HsExtendRightState = PlayerState{
	Loop: func(p *Player) {
		p.HsExtendLoop(direction.MaxRight())
	},
}.denil()

var HsExtendLeftState = PlayerState{
	Loop: func(p *Player) {
		p.HsExtendLoop(direction.MaxLeft())
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
		//p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
		p.HsPullLoop(direction.MaxRight())
		//p.PullPlayer()
	},
}.denil()

var HsPullLeftState = PlayerState{
	Loop: func(p *Player) {
		p.HsPullLoop(direction.MaxLeft())
		//p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
		//p.PullPlayer()
	},
}.denil()

var HsPullUpState = PlayerState{
	Loop: func(p *Player) {
		p.HsPullLoop(direction.MaxUp())
	},
}

var HsItemGrabRightState = PlayerState{
	Start: func(p *Player) {
		p.HeldObj = p.Hs.GetLastHitObj(true)
	},
	Loop: func(p *Player) {
		p.HsItemGrabLoop(direction.MaxRight())
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
		p.HsItemGrabLoop(direction.MaxLeft())
	},
	End: func(p *Player) {
		p.HeldObj.Delta.SetPos(0,0)
		p.HeldObj = nil
	},
}.denil()
