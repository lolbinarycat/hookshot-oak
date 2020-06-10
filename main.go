package main

import (
	"image/color"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

const Ground collision.Label = 1

type ControlConfig struct {
	Left, Right, Jump string
}

//type Player struct {
//
//	Body
//}

var currentControls ControlConfig = ControlConfig{
	Left:  key.LeftArrow,
	Right: key.RightArrow,
	Jump:  key.Z,
}

func main() {
	oak.Add("platformer", func(string, interface{}) {
		playerBody := entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)

		render.Draw(playerBody.R)
		playerBody.Speed = physics.NewVector(3, 3)

		ground := entities.NewSolid(0, 400, 500, 20,
			render.NewColorBox(500, 20, color.RGBA{0, 0, 255, 255}),
			nil, 0)
		ground.UpdateLabel(Ground)

		render.Draw(ground.R)

		playerBody.Bind(func(id int, nothing interface{}) int {
			if oak.IsDown(currentControls.Left) {
				playerBody.ShiftX(-playerBody.Speed.X())
			}
			if oak.IsDown(currentControls.Right) {
				playerBody.ShiftX(playerBody.Speed.X())
			}

			//gravity
			fallSpeed := .1
			onGround := playerBody.HitLabel(Ground)
			if onGround == nil {
				playerBody.Delta.ShiftY(fallSpeed)
			} else {
				playerBody.Delta.SetY(0)

				if oak.IsDown(currentControls.Jump) {
					playerBody.Delta.ShiftY(-playerBody.Speed.Y())
				}
			}
			//apply gravity
			playerBody.ShiftY(playerBody.Delta.Y())

			return 0
		}, event.Enter)
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})

	oak.Init("platformer")
}
