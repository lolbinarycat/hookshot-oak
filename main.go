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
	Left, Right, Jump, Jump2, UpwardForce string
}

var currentControls ControlConfig = ControlConfig{
	Left:        key.LeftArrow,
	Right:       key.RightArrow,
	Jump:        key.Z,
	Jump2:       key.X,
	UpwardForce: key.A,
}

//type State struct {
//}

var player Player

type Player struct {
	Body  *entities.Moving
	State func() func()
}

func (p Player) AirState() func() { //start in air state
	//gravity
	fallSpeed := .1

	if isOnGround(p.Body) {
		return func() { player.State = player.GroundState }
		//panic("t")
	} else {
		p.Body.Delta.ShiftY(fallSpeed)

	}
	p.Body.ShiftY(p.Body.Delta.Y())

	//return p.State
	//panic("e")
	return func() {}
}

func (p Player) GroundState() func() {
	if isOnGround(p.Body) {

		p.Body.Delta.SetY(0)

		if oak.IsDown(currentControls.Jump2) {

			p.Body.Delta.ShiftY(-p.Body.Speed.Y())
			p.Body.ShiftY(p.Body.Delta.Y())
			return func() { player.State = player.AirState }
		}
	} else {
		return func() { player.State = player.AirState }
	}
	return func() { player.State = player.AirState }
}

func isOnGround(mov *entities.Moving) bool {
	onGround := mov.HitLabel(Ground)
	if onGround == nil {
		return false
	} else {
		return true
	}
}

func main() {
	oak.Add("platformer", func(string, interface{}) {

		player.Body = entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)
		player.State = player.AirState

		render.Draw(player.Body.R)
		player.Body.Speed = physics.NewVector(3, 3)

		ground := entities.NewSolid(0, 400, 500, 20,
			render.NewColorBox(500, 20, color.RGBA{0, 0, 255, 255}),
			nil, 0)
		ground.UpdateLabel(Ground)

		render.Draw(ground.R)

		player.Body.Bind(func(id int, nothing interface{}) int {

			if oak.IsDown(currentControls.Left) {
				player.Body.ShiftX(-player.Body.Speed.X())
			}
			if oak.IsDown(currentControls.Right) {
				player.Body.ShiftX(player.Body.Speed.X())
			}

			player.State()()

			return 0
		}, event.Enter)
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})

	oak.Init("platformer")
}
