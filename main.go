package main

import (
	"image/color"
	"os"
	"time"

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

const JumpHeight int = 3
const CoyoteTime time.Duration = time.Millisecond * 7

//type CollisionType int8

type ActiveCollisions struct {
	GroundHit    bool
	LeftWallHit  bool
	RightWallHit bool
	CeilingHit   bool
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
	Body *entities.Moving
	ActiColls ActiveCollisions
	State      func()
	StateStartTime time.Time
	//StateInit  bool
}

func (p Player) AirState() { //start in air state
	//gravit
	fallSpeed := .1
	//hitType, _, _ := howIsHittingLabel(p.Body, Ground)
	if player.ActiColls.GroundHit {
		//hit := collision.HitLabel(p.Body.Space, Ground)
		// Correct our y if we started falling into the ground
		//p.Body.SetY(hit.Y() - p.Body.H)
		//p.Body.Delta.SetY(0)
		print("g")
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
	//hitType,_,_ := howIsHittingLabel(p.Body,Ground)
	if player.ActiColls.GroundHit {
		p.Body.Delta.SetY(0)

		if oak.IsDown(currentControls.Jump) {
			p.Jump()
		} //else {
		//p.Body.Delta.ShiftY(fallSpeed)
		//}
	} else {
		print("c")
		p.SetState(player.CoyoteState)
	}
	//howIsHittingLabel(p.Body, Ground)
}

//CoyoteState implements "coyote time" a window of time after
//running off an edge in which you can still jump
func (p Player) CoyoteState() {
	if p.StateStartTime.Add(CoyoteTime).Before( time.Now()) {
		p.SetState(p.AirState)
	} 
	//inherit code from AirState
	p.AirState()
	if oak.IsDown(currentControls.Jump) {
		p.Jump()
	}

}

func (p Player) Jump() {
	p.Body.Delta.ShiftY(-p.Body.Speed.Y())
	p.Body.ShiftY(p.Body.Delta.Y())
	player.SetState(p.AirState)
}

func (p Player) SetState(state func()) {
	player.StateStartTime = time.Now()
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

func isInGround(mov *entities.Moving) bool {
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

	ground := entities.NewSolid(0, 400, 500, 20,
		render.NewColorBox(500, 20, color.RGBA{0, 0, 255, 255}),
		nil, 0)
	ground2 := entities.NewSolid(0, 200, 20, 500,
		render.NewColorBox(20, 500, color.RGBA{0, 255, 255, 255}),
		nil, 1)
	ground.UpdateLabel(Ground)
	ground2.UpdateLabel(Ground)

	render.Draw(ground.R)
	render.Draw(ground2.R, 1)

}

func main() {
	oak.Add("platformer", func(string, interface{}) {
		loadScene()

		player.Body.Bind(func(id int, nothing interface{}) int {

			if oak.IsDown(currentControls.Left) {
				player.Body.Delta.SetX(-player.Body.Speed.X())
			} else if oak.IsDown(currentControls.Right) {
				player.Body.Delta.SetX(player.Body.Speed.X())
			} else {
				player.Body.Delta.SetX(0)
			}
			if oak.IsDown(currentControls.Quit) {
				os.Exit(0)
			}

			//HowIsHittingLabel(player.Body,Ground)
			//this is all it takes to make collision, why is everyone overcomplicating it? <- this statement won't age well
			oldX, oldY := player.Body.GetPos()
			player.State()
			player.ActiColls = ActiveCollisions{}
			player.Body.ShiftX(player.Body.Delta.X())
			if collision.HitLabel(player.Body.Space, Ground) != nil {
				player.Body.SetX(oldX)
			}

			player.Body.ShiftY(player.Body.Delta.Y())

			if collision.HitLabel(player.Body.Space, Ground) != nil {
				if player.Body.Delta.Y() > 0 {
					player.ActiColls.GroundHit = true
				} else if player.Body.Delta.Y() < 0 {
					player.ActiColls.CeilingHit = true
				}
				
				//  player.Body.Delta.SetY(0)
				player.Body.SetY(oldY)
			
			}

			return 0
		}, event.Enter)
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})

	oak.Init("platformer")
}
