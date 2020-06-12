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
	LeftBottomCornerHit = 3
)
 
func howIsHittingLabel(mov *entities.Moving, label collision.Label) (cType CollisionType,xOverlap ,yOverlap float64) {
	//oldX, _ := mov.GetPos()
	hitArr := collision.Hits(mov.Space)
	if len(hitArr) == 2 {
		return LeftBottomCornerHit, 0,0
	} 
	
	if len(hitArr)== 1 {
		hit := hitArr[0]
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
			/*xbump := 1.0
			if xover > 0 {
				xbump = -1
			}
			mov.SetX(oldX + xbump)
			if mov.Delta.Y() < 0 {
				mov.Delta.SetY(0)
			}*/
			//	print("left")
			return LeftWallHit,xover,yover
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
	return GroundHit, 0, 0
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
	hitType,_,_ := howIsHittingLabel(p.Body,Ground)
	if hitType == GroundHit {
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
	//hitType,_,_ := howIsHittingLabel(p.Body,Ground)
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
		p.StateTimer = 6
		p.StateInit = false
	} else if p.StateTimer <= 0 {
		p.SetState(p.AirState)
	} else {
		p.StateTimer--
	}
	/*if (isInGround(p.Body)) {
		hit := collision.HitLabel(p.Body.Space, Ground)
		// Correct our y if we started falling into the ground
		p.Body.SetY(hit.Y() - p.Body.H)
		p.Body.Delta.SetY(0)
		player.State = player.GroundState*/
	 if isOnGround(p.Body) {
		player.State = player.GroundState
		//return func() {}
		//panic("t")

	} else {
		//p.SetState(p.AirState)
		p.Body.Delta.ShiftY(0.1)

	}
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
			//this is all it takes to make collision, why is everyone overcomplicating it?
			oldX, oldY := player.Body.GetPos()
			player.State()
			player.Body.ShiftY(player.Body.Delta.Y())
			if collision.HitLabel(player.Body.Space,Ground) != nil {
				player.Body.Delta.SetY(0)
				player.Body.SetY(oldY)
			}
			player.Body.ShiftX(player.Body.Delta.X())
			if collision.HitLabel(player.Body.Space,Ground) != nil {
				player.Body.SetX(oldX)
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
