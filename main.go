package main

import (
	"image/color"
	"math"

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

const JumpHeight int = 5

type CollisionType int16

const (
	GroundHit    CollisionType = 0
	LeftWallHit                = 1
	RightWallHit               = 2
)

func HowIsHittingLabel(mov *entities.Moving) {
	oldX, _ := mov.GetPos()
	hit := collision.HitLabel(mov.Space, Ground)
	if hit != nil {
		// If we walked into a piece of ground, move back
		xover, yover := mov.Space.Overlap(hit)
		// We, perhaps unintuitively, need to check the Y overlap, not
		// the x overlap
		// if the y overlap exceeds a superficial value, that suggests
		// we're in a state like
		//
		// G = Ground, C = Movacter
		//
		// GG C
		// GG C
		//
		// moving to the left
		if math.Abs(yover) > 1 {
			// We add a buffer so this doesn't retrigger immediately
			xbump := 1.0
			if xover > 0 {
				xbump = -1
			}
			mov.SetX(oldX + xbump)
			if mov.Delta.Y() < 0 {
				mov.Delta.SetY(0)
			}
		}

		// If we're below what we hit and we have significant xoverlap, by contrast,
		// then we're about to jump from below into the ground, and we
		// should stop the movacter.
		//if !aboveGround && math.Abs(xover) > 1 {
		//	// We add a buffer so this doesn't retrigger immediately
		//	mov.SetY(oldY + 1)
		//	mov.Delta.SetY(fallSpeed)
		//}
	}

}

type ControlConfig struct {
	Left, Right, Jump string
}

var currentControls ControlConfig = ControlConfig{
	Left:  key.LeftArrow,
	Right: key.RightArrow,
	Jump:  key.Z,
	//Jump2:       key.X,
	//	UpwardForce: key.A,
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
		player.State = player.GroundState
		return func() {}
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

		if oak.IsDown(currentControls.Jump) {

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

func loadScene() {


		player.Body = entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)
		player.State = player.AirState

		render.Draw(player.Body.R)
		player.Body.Speed = physics.NewVector(3, float64(JumpHeight))

		ground := entities.NewSolid(0, 400, 500, 20,
			render.NewColorBox(500, 20, color.RGBA{0, 0, 255, 255}),
			nil, 0)
		ground2 := entities.NewSolid(0,400 , 500, 20,
			render.NewColorBox(500,20, color.RGBA{0, 0, 255, 255}),
			nil, 0)
		ground.UpdateLabel(Ground)
		ground2.UpdateLabel(Ground)

		render.Draw(ground.R)
		render.Draw(ground.R)
	
}

func main() {
	oak.Add("platformer", func(string, interface{}) {
	loadScene()

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
