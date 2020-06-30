package main

import (
	"time"

	oak "github.com/oakmound/oak/v2"
)

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
			p.SetState(GroundPoundEndState)
		} else if oak.IsDown(curCtrls.Up){
			p.Body.Delta.SetY(0)
			p.SetState(AirState)
		} else {
			p.Body.Delta.SetY(GroundPoundSpeed)
		}
	},
}.denil()

const GroundPoundEndTime = time.Millisecond * 80
var GroundPoundEndState = PlayerState{
	Loop: func(p *Player) {
		if p.TimeFromStateStart() > GroundPoundEndTime {
			p.SetState(GroundState)
		}
		p.Body.Delta.SetX(0)
	},
}.denil()
