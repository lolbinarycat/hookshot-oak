package main

import (
	"fmt"
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

const JumpHeight int = 6
const Gravity  float64 = 0.2
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
	//Body           *entities.Moving
	//ActiColls      ActiveCollisions
	PhysObject
	State          func()
	StateStartTime time.Time
}

type PhysObject struct {
	Body *entities.Moving
	ActiColls ActiveCollisions
}

func (p *Player) AirState() { //start in air state
	//print("a")
	fallSpeed := .1

	if player.PhysObject.ActiColls.GroundHit {
		player.State = player.GroundState
		//print("ground")
	} else {
		p.Body.Delta.ShiftY(fallSpeed)

	}
	//p.Body.ShiftY(p.Body.Delta.Y())

	//return p.State
	//panic("e")
}

func (p *Player) GroundState() {
	//fallSpeed := .1
	//print("groundstate")
	//hitType,_,_ := howIsHittingLabel(p.Body,Ground)
	if player.PhysObject.ActiColls.GroundHit == true {
		//p.Body.Delta.SetY(0)

		if oak.IsDown(currentControls.Jump) {
			p.Jump()
		}
	} else {
		//print("c")
		p.SetState(player.CoyoteState)
	}
	//p.Body.Delta.ShiftY(0.001)
	//howIsHittingLabel(p.Body, Ground)
}

//CoyoteState implements "coyote time" a window of time after
//running off an edge in which you can still jump
func (p Player) CoyoteState() {
	if p.StateStartTime.Add(CoyoteTime).Before(time.Now()) {
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

func (p *Player) SetState(state func()) {
	p.StateStartTime = time.Now()
	p.State = state
}

//TimeFromStateStart gets how long it has been since the last state transition
func (p Player) TimeFromStateStart() time.Duration {
	return p.StateStartTime.Sub(time.Now())
}

//this functon should not be used, use Body.ActiColls.GroundHit instead
func isOnGround(mov *entities.Moving) bool {
	onGround := mov.HitLabel(Ground)
	if onGround == nil {
		return false
	} else {
		return true
	}
}

// Depreciated, do not use
func isInGround(mov *entities.Moving) bool {
	_, oldY := mov.GetPos()
	hit := collision.HitLabel(mov.Space, Ground)
	if hit != nil && !(oldY != mov.Y() && oldY+mov.H > hit.Y()) {
		return true
	}

	return false
}

func (object *PhysObject) doCollision(updater func()) {
	oldX, oldY := object.Body.GetPos()
	updater()
	object.ActiColls = ActiveCollisions{} //reset the struct to be all false
	

	object.Body.ShiftX(object.Body.Delta.X())
	if hit := collision.HitLabel(object.Body.Space, Ground); hit != nil {
		//xover, _ := object.Body.Space.Overlap(hit)
		if object.Body.Delta.X() > 0 {
			object.ActiColls.RightWallHit = true
			object.Body.SetX(oldX)
		} else if object.Body.Delta.X() < 0 {
			object.ActiColls.LeftWallHit = true
			//object.Body.SetX(object.Body.X() - xover)
			//object.Body.SetX(oldX)
			object.Body.SetX(hit.X()+hit.W())
		}
	}

	object.Body.ShiftY(object.Body.Delta.Y())
	if hit := collision.HitLabel(object.Body.Space, Ground); hit != nil {
		//_, yover := object.Body.Space.Overlap(hit)
		if object.Body.Delta.Y() > 0 {
			object.ActiColls.GroundHit = true
			//object.Body.SetY(object.Body.Y() - yover)
			//object.Body.Delta.SetY(0.7*yover)
			//print("u")
			//object.Body.SetY(object.Body.Y()-object.Body.Delta.Y())
			object.Body.SetY(hit.Y()-object.Body.H)
		} else if object.Body.Delta.Y() < 0 {
			object.ActiColls.CeilingHit = true
			object.Body.SetY(oldY)
		}
		object.Body.Delta.SetY(0)
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
	ground2 := entities.NewSolid(0, 200, 20, 500,
		render.NewColorBox(20, 500, color.RGBA{0, 255, 255, 255}),
		nil, 1)
	ground3 :=  entities.NewSolid(300, 200, 20, 500,
		render.NewColorBox(20, 500, color.RGBA{0, 255, 255, 255}),
		nil, 2)

	blockArr := make([]*entities.Solid,10)
	for i := 0;i < 10;i++ {
		blockArr[i] = entities.NewSolid(float64(10*i), 50,10,10, render.NewColorBox(10,10,color.RGBA{255,0,uint8(16*i),255}), nil, event.CID(i+3))
		render.Draw(blockArr[i].R,8+i)
	}
	ground.UpdateLabel(Ground)
	ground2.UpdateLabel(Ground)
	ground3.UpdateLabel(Ground)

	render.Draw(ground.R)
	render.Draw(ground2.R, 1)
	render.Draw(ground3.R, 1)

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
				if oak.IsDown(key.I) {
					fmt.Println(player)
				}
				os.Exit(0)
			}


			player.doCollision(player.State)

			return 0
		}, event.Enter)
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})

	oak.Init("platformer")
}


