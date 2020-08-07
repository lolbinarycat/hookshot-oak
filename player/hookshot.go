package player

import (
	"time"

	"github.com/lolbinarycat/hookshot-oak/direction"
	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/lolbinarycat/hookshot-oak/physobj"
	"github.com/oakmound/oak/v2/dlog"
)

func (p *Player) EndHs() {
	p.Hs.Active = false
	p.Hs.X = 0
	p.Hs.Y = 0
	p.Hs.Body.Delta.SetPos(0, 0)
}

func GiveMods(mods ...*PlayerModule) {
	for _, m := range mods {
		(*m).Obtain()
		if autoEquipMods {
			(*m).Equip()
		}
	}
}

func (p *Player) HsUpdater() {
	hsOffX := p.Body.W/2 - p.Hs.Body.H/2
	hsOffY := p.Body.H/2 - p.Hs.Body.H/2

	//set hookshot's relitive position to be accurate
	p.Hs.X = p.Hs.Body.X() - p.Body.X() - hsOffX
	p.Hs.Y = p.Hs.Body.Y() - p.Body.Y() - hsOffY
}

// This file contains functions that handle hookshot behavior

const HsStartTime time.Duration = time.Millisecond * 60

var HsStartState = PlayerState{
	Loop: func(p *Player) {
		if p.Mods["hs"].Active() == false {
			p.SetState(AirState)
			return
		}
		if p.TimeFromStateStart() > HsStartTime {
			dir := p.HeldDir
			if (dir == direction.Dir{}) {
				// if no direction is held, return to airstate
				p.SetState(AirState)
			} else {
				if p.Mods["diaghs"].Active() {
					dir = dir.OrthoDiagonalize()
				} else {
					dir = dir.Orthogonalize()
				}
				p.SetState(HsExtendState(dir.Maximize()))
			}
		}
	}}

const HsExtendTime time.Duration = time.Second * 2

func HsExtendState(dir direction.Dir) PlayerState {
	dlog.Verb("hookshot extending with dir",dir)
	coeffX := direction.ToCoeff(dir.H)
	coeffY := direction.ToCoeff(dir.V)
	return PlayerState{
		Loop: func(p *Player) {
			p.Hs.Active = true
			if p.TimeFromStateStart() > HsExtendTime {
				goto Retract
			} else if p.TimeFromStateStart() > HsInputTime && p.IsHsInput() {
				goto Retract
			} else {
				if dir.IsRight() && p.Hs.ActiColls.RightWallHit ||
					dir.IsLeft() && p.Hs.ActiColls.LeftWallHit ||
					dir.IsUp() && p.Hs.ActiColls.CeilingHit ||
					dir.IsDown() && p.Hs.ActiColls.GroundHit {

					if (dir.H == 0 || p.Hs.ActiColls.HLabel == labels.NoHs) &&
					   (dir.V == 0 || p.Hs.ActiColls.VLabel == labels.NoHs) {
						goto Retract
					}

					if p.Mods["hsitemgrab"].Active() &&
						((p.Hs.ActiColls.HLabel == labels.Block && dir.H != 0) ||
						(p.Hs.ActiColls.VLabel == labels.Block && dir.V != 0)) {
						p.SetState(HsItemGrabState(dir))
					} else {
						p.SetState(HsPullState(dir))
					}
				}
				p.Body.Delta.SetPos(0, 0)
				p.Hs.Body.Delta.SetPos(p.Hs.Body.Speed.X()*coeffX,
					p.Hs.Body.Speed.Y()*coeffY)

				}
				return
			Retract:
				p.SetState(HsRetractState(dir))

		},
	}
}

func HsRetractState(dir direction.Dir) PlayerState {
	coeffX := direction.ToCoeff(dir.H)
	coeffY := direction.ToCoeff(dir.V)
	return PlayerState{
		Loop: func(p *Player) {
			if (dir.IsLeft() && p.Hs.X >= 0) ||
				(dir.IsRight() && p.Hs.X <= 0) ||
				(dir.IsUp() && p.Hs.Y >= 0) ||
				(dir.IsDown() && p.Hs.Y <= 0) {
				p.EndHs()
				p.SetState(AirState)
				return
			}

			p.Hs.Body.Delta.SetPos(-p.Hs.Body.Speed.X()*coeffX,
				-p.Hs.Body.Speed.Y()*coeffY)
		},
	}
}

func HsPullState(dir direction.Dir) PlayerState {
	coeffX := direction.ToCoeff(dir.H)
	coeffY := direction.ToCoeff(dir.V)
	return PlayerState{
		Start: func(p *Player) {
			dlog.Verb("hsActive =",p.Hs.Active)
		},
		Loop: func(p *Player) {
			if p.HasHitInDir(dir) {
				p.EndHs()
				p.SetState(AirState)
				return
			}
			p.Body.Delta.SetPos(coeffX*p.Hs.Body.Speed.X(),
				coeffY*p.Hs.Body.Speed.Y())
			p.Hs.Delta.SetPos(0,0)
		},
	}
}

func HsItemGrabState(dir direction.Dir) PlayerState {
	coeffX := direction.ToCoeff(dir.H)
	coeffY := direction.ToCoeff(dir.V)
	return PlayerState{
		Start: func(p *Player) {
				var ok bool
				p.HeldObj, ok = p.Hs.ActiColls.LastHitH.E().(*physobj.Block)
				if !ok {
					p.HeldObj = p.Hs.ActiColls.LastHitV.E().(*physobj.Block)
				}
		},
		Loop: func(p *Player) {
			if (dir.IsJustRight() && (p.Hs.X <= 0 || p.ActiColls.RightWallHit)) ||
				(dir.IsJustLeft() && (p.Hs.X >= 0 || p.ActiColls.LeftWallHit)) ||
				(dir.IsUp() && (p.Hs.Y >= 0 || p.ActiColls.CeilingHit)) ||
				(dir.IsDown() && (p.Hs.Y <= 0 || p.ActiColls.GroundHit)) {
				p.EndHs()
				p.State.End(p)

				return
			}
			p.Hs.Body.Delta.SetPos(-p.Hs.Body.Speed.X()*coeffX,
				-p.Hs.Body.Speed.X()*coeffY)
			p.HeldObj.Body.Delta.SetPos(p.Hs.Body.Delta.GetPos())
		},
		End: func(p *Player) {
			if p.Mods["itemcarry"].Active() {
				p.SetStateAdv(ItemCarryGroundState, SetStateOptArgs{SkipEnd: true})
			} else {
				p.HeldObj.Body.Delta.SetPos(0, 0)
				p.HeldObj = nil
				p.SetStateAdv(AirState,SetStateOptArgs{SkipEnd: true})
			}
		},
	}
}
