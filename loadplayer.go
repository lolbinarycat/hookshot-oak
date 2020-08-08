package main

import (
	"image/color"

	"github.com/lolbinarycat/hookshot-oak/labels"
	//"github.com/lolbinarycat/hookshot-oak/layers"
	"github.com/lolbinarycat/hookshot-oak/player"
	prender "github.com/lolbinarycat/hookshot-oak/player/renderable"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/physics"
	"github.com/oakmound/oak/v2/render"
)

const (
	PlayerWidth = prender.PlayerWidth
	PlayerHeight = prender.PlayerHeight
)

const HsWidth = 4
const HsHeight = 4


func loadPlayer() *player.Player {
	var err error
	var plr = new(player.Player)
	plr.R, err = prender.LoadCom()
	if err != nil {
		panic(err)
	}

	plr.Body = entities.NewMoving(300, 400, PlayerWidth, PlayerHeight,
		plr.R,
		nil, 0, 0)
	plr.Body.Init()
	plr.Space.UpdateLabel(labels.Player)

	player.SetPlayer(0, plr)

	plr.State = &player.RespawnFallState
	plr.RespawnPos = player.Pos{X: plr.Body.X(), Y: plr.Body.Y()}
	render.Draw(plr.Body.R, 1)
	plr.Body.Speed = physics.NewVector(3, float64(player.JumpHeight))

	plr.Hs.Body = &*entities.NewMoving(100, 100, HsWidth, HsHeight,
		render.NewColorBox(HsHeight, HsWidth, color.RGBA{0, 0, 255, 255}),
		nil, 1, 0)
	plr.Hs.Body.Init()

	return plr
}

