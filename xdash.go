package main

import (
	"time"

	"github.com/oakmound/oak/v2"
)

const XDashSpeed = 8

var XDashState = PlayerState{
	Loop: func(p *Player) {
		if oak.IsDown(p.Ctrls.Left){
			p.Body.Delta.SetX(-XDashSpeed)
		} else if oak.IsDown(p.Ctrls.Right) {
			p.Body.Delta.SetX(XDashSpeed)
		}
	
			if isJumpInput() {
				p.Jump()
			}

	
	},
	MaxDuration: time.Millisecond * 400,
}.denil()
