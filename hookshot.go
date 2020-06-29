package main

import (
	"time"

	"github.com/lolbinarycat/hookshot-oak/direction"
	"github.com/lolbinarycat/hookshot-oak/labels"

	oak "github.com/oakmound/oak/v2"
)

// This file contains functions that handle hookshot behavior

const HsStartTime time.Duration = time.Millisecond * 60

var HsStartState = PlayerState{
	Loop: func(p *Player) {
		if player.Mods.Hookshot.Equipped == false {
			p.SetState(AirState)
			return
		}
		if p.TimeFromStateStart() > HsStartTime {
			if oak.IsDown(currentControls.Right) {
				p.SetState(HsExtendState(direction.MaxRight()))
			} else if oak.IsDown(currentControls.Left) {
				p.SetState(HsExtendState(direction.MaxLeft()))
			} else {
				p.SetState(AirState)
			}
		}
	}}.denil()

func HsExtendState(dir direction.Dir) PlayerState {
	return PlayerState{
		Loop: func(p *Player) {
			p.Hs.Active = true

			if p.TimeFromStateStart() > HsExtendTime {
				p.SetState(HsRetractState(dir))
			} else if (dir.IsRight() && p.Hs.ActiColls.RightWallHit) ||
				(dir.IsLeft() && p.Hs.ActiColls.LeftWallHit) {
				if p.Hs.ActiColls.HLabel == labels.Block &&
					p.Mods.HsItemGrab.Equipped {
					p.SetState(HsItemGrabState(dir))
				} else {
					p.SetState(HsPullState(dir))
				}
			} else if p.TimeFromStateStart() > HsInputTime && isHsInput() {
				p.SetState(HsRetractState(dir))
			} else {
				p.Body.Delta.SetPos(0, 0)
				if dir.IsJustRight() {
					p.Hs.Body.Delta.SetX(p.Hs.Body.Speed.X())
				} else if dir.IsJustLeft() {
					p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
				}
			}
		},
	}.denil()
}

func HsRetractState(dir direction.Dir) PlayerState {
	coeffX := direction.ToCoeff(dir.H)
	coeffY := direction.ToCoeff(dir.V)
	return PlayerState{
		Loop: func(p *Player) {
			if (dir.IsLeft() && p.Hs.X >= 0) ||
				(dir.IsRight() && p.Hs.X <= 0) {
				p.EndHs()
				return
			}

			p.Hs.Body.Delta.SetPos(-p.Hs.Body.Speed.X()*coeffX,
				-p.Hs.Body.Speed.Y()*coeffY)
		},
	}.denil()
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
			p.Body.Delta.SetPos(coeffX*p.Hs.Body.Speed.X(),
				coeffY*p.Hs.Body.Speed.Y())
		},
	}.denil()
}

func HsItemGrabState(dir direction.Dir) PlayerState {
	coeffX := direction.ToCoeff(dir.H)
	coeffY := direction.ToCoeff(dir.V)
	return PlayerState{
		Start: func(p *Player) {
			p.HeldObj = p.Hs.GetLastHitObj(true)
		},
		Loop: func(p *Player) {
			if (dir.IsJustRight() && (p.Hs.X <= 0 || p.ActiColls.RightWallHit)) ||
				(dir.IsJustLeft() && (p.Hs.X >= 0 || p.ActiColls.LeftWallHit)) {
				p.EndHs()
				return
			}
			p.Hs.Body.Delta.SetPos(-p.Hs.Body.Speed.X()*coeffX,
				-p.Hs.Body.Speed.X()*coeffY)
			p.HeldObj.Delta.SetPos(p.Hs.Body.Delta.GetPos())
		},
		End: func(p *Player) {
			p.HeldObj.Delta.SetPos(0, 0)
			p.HeldObj = nil
		},
	}.denil()
}

//STATES:

// var HsRetractRightState = PlayerState{
// 	Loop: func(p *Player) {
// 		if p.Hs.X <= 0 {
// 			p.EndHs()
// 			return
// 		}
// 		p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
// 		//p.Hs.X -= p.Hs.Body.Speed.X()
// 	},
// }.denil()

// var HsRetractLeftState = PlayerState{
// 	Loop: func(p *Player)
// }.denil()

//HsPullRightState is the state for when the hookshot is
//pulling the player after having hit an object
// var HsPullRightState = PlayerState{
// 	Loop: func(p *Player) {
// 		//p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
// 		p.HsPullLoop(direction.MaxRight())
// 		//p.PullPlayer()
// 	},
// }.denil()

// var HsPullLeftState = PlayerState{
// 	Loop: func(p *Player) {
// 		p.HsPullLoop(direction.MaxLeft())
// 		//p.Hs.Body.Delta.SetX(-p.Hs.Body.Speed.X())
// 		//p.PullPlayer()
// 	},
// }.denil()

//var HsPullUpState = PlayerState{
//	Loop: func(p *Player) {
