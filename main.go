package main

import (
	"image/color"
	"math"
	"os"

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

type CollisionType int8

const (
	GroundHit    CollisionType = 0
	LeftWallHit                = 1
	RightWallHit               = 2
)
 
func howIsHittingLabel(mov *entities.Moving, label collision.Label) CollisionType{
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
			xbump := 1.0
			if xover > 0 {
				xbump = -1
			}
			mov.SetX(oldX + xbump)
			if mov.Delta.Y() < 0 {
				mov.Delta.SetY(0)
			}
			return LeftWallHit
			// We add a buffer so this doesn't retrigger immediately
			///
			//*/
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
	return GroundHit
}

type ControlConfig struct {
	Left, Right, Jump, Quit string
}

var currentControls ControlConfig = ControlConfig{
	Left:  key.LeftArrow,
	Right: key.RightArrow,
	Jump:  key.Z,
	Quit:  key.Q,
	//Jump2:       key.X,
	//	UpwardForce: key.A,
}

//type State struct {
//}

var player Player

//Player is a type representing the player
//StateInit is a variable that should be set to true when changing states
//it tells the state to initialize values like StateTimer
type Player struct {
	Body  *entities.Moving
	State func()
	StateTimer int64
	StateInit bool
}

func (p Player) AirState()  { //start in air state
	//gravit
	fallSpeed := .1

	if (howIsHittingLabel(p.Body,Ground) == GroundHit) {
		//hit := collision.HitLabel(p.Body.Space, Ground)
		// Correct our y if we started falling into the ground
		//p.Body.SetY(hit.Y() - p.Body.H)
		//p.Body.Delta.SetY(0)
		player.State = player.GroundState
		//print("ground")
	} else if isOnGround(p.Body) {
		 player.State = player.GroundState
		 //return func() {}
		 //panic("t")

	 } else {
		p.Body.Delta.ShiftY(fallSpeed)

	}
	p.Body.ShiftY(p.Body.Delta.Y())

	//return p.State
	//panic("e")
}

func (p Player) GroundState() {
	//fallSpeed := .1
	//print("groundstate")
	if isOnGround(p.Body) {

		p.Body.Delta.SetY(0)

		if oak.IsDown(currentControls.Jump) {
			p.Jump()
		} //else {
		//p.Body.Delta.ShiftY(fallSpeed)
		//}
	} else {
		//print("air")
		p.SetState(player.CoyoteState)
	}
	howIsHittingLabel(p.Body,Ground)
}


func (p Player) CoyoteState() {
	if p.StateInit {
		p.StateTimer = 10
		p.StateInit = false
	} else if p.StateTimer <= 0 {
		p.SetState(p.AirState)
	} else {
		p.StateTimer--
	}
	if (isInGround(p.Body)) {
		hit := collision.HitLabel(p.Body.Space, Ground)
		// Correct our y if we started falling into the ground
		p.Body.SetY(hit.Y() - p.Body.H)
		p.Body.Delta.SetY(0)
		player.State = player.GroundState
		print("ground")
	} else if isOnGround(p.Body) {
		player.State = player.GroundState
		//return func() {}
		//panic("t")

	} else {
		//p.SetState(p.AirState)
		p.Body.Delta.ShiftY(0.1)

	}
	p.Body.ShiftY(p.Body.Delta.Y())
	if oak.IsDown(currentControls.Jump) {
		p.Jump()
	}

}

func (p Player) Jump() {
	p.Body.Delta.ShiftY(-p.Body.Speed.Y())
	p.Body.ShiftY(p.Body.Delta.Y())
	p.SetState(p.AirState)
}

func (p Player) SetState(state func()) {
	player.StateInit = true
	player.State = state
}

func isOnGround(mov *entities.Moving) bool {
	onGround := mov.HitLabel(Ground)
	if onGround == nil {
		return false
	} else {
		return true
	}
}

func isInGround (mov *entities.Moving) bool {
	_, oldY := mov.GetPos()
	hit := collision.HitLabel(mov.Space, Ground)
	if hit != nil && !(oldY != mov.Y() && oldY+mov.H > hit.Y()) {
		return true
	}

	return false
}

func loadScene() {

		player.Body = entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)
		player.State = player.AirState

		render.Draw(player.Body.R)
		player.Body.Speed = physics.NewVector(3, float64(JumpHeight))

		ground  := entities.NewSolid(0, 400, 500, 20,
			render.NewColorBox(500, 20, color.RGBA{0, 0, 255, 255}),
			nil, 0)
		ground2 := entities.NewSolid(0, 200, 20, 500,
			render.NewColorBox(20,500, color.RGBA{0, 255, 255, 255}),
			nil, 1)
		ground.UpdateLabel(Ground)
		ground2.UpdateLabel(Ground)

		render.Draw(ground.R)
		render.Draw(ground2.R,1)
	
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
			if oak.IsDown(currentControls.Quit) {
				os.Exit(0)
			}

			//HowIsHittingLabel(player.Body,Ground)
			player.State()

			return 0
		}, event.Enter)
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})

	oak.Init("platformer")
}
