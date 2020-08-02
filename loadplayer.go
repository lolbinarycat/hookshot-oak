package main

import (
	"image/color"

	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/lolbinarycat/hookshot-oak/layers"
	"github.com/lolbinarycat/hookshot-oak/player"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/physics"
	"github.com/oakmound/oak/v2/render"
)

const PlayerWidth = 12
const PlayerHeight = 12

const HsWidth = 4
const HsHeight = 4


func loadPlayer() *player.Player {
	var eyeColor = color.RGBA{0, 255, 255, 255}
	playerSprite, err := render.LoadSprite("assets/images", "player_new.png")
	if err != nil {
		panic(err)
	}
	var plr = new(player.Player)
	eye1 := render.NewColorBox(1, 4, eyeColor)
	eye2 := eye1.Copy().(*render.Sprite)
	plr.Eyes = [2]*render.Sprite{eye1, eye2}
	plr.Body = entities.NewMoving(300, 400, PlayerWidth, PlayerHeight,
		playerSprite,
		nil, 0, 0)
	plr.Body.Init()
	plr.Space.UpdateLabel(labels.Player)

	eye1.LayeredPoint.Vector = eye1.Attach(plr.Body, 4, 3)
	eye2.LayeredPoint.Vector = eye1.Attach(plr.Body, 8, 3)

	player.SetPlayer(0, plr)

	render.Draw(eye1, layers.FG, 2)
	render.Draw(eye2, layers.FG, 2)

	plr.State = player.RespawnFallState
	plr.RespawnPos = player.Pos{X: plr.Body.X(), Y: plr.Body.Y()}
	render.Draw(plr.Body.R, 1)
	plr.Body.Speed = physics.NewVector(3, float64(player.JumpHeight))

	plr.Hs.Body = &*entities.NewMoving(100, 100, HsWidth, HsHeight,
		render.NewColorBox(HsHeight, HsWidth, color.RGBA{0, 0, 255, 255}),
		nil, 1, 0)
	plr.Hs.Body.Init()

	return plr
}

