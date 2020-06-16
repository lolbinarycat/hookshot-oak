package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"time"

	//"compress/flate"

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
const WallJumpHeight float64 = 6
const WallJumpWidth float64 = 3
const WallJumpLaunchDuration time.Duration = time.Millisecond * 230
const (
	AirAccel float64 = 0.4
	AirMaxSpeed float64 = 3
)
//Window constants
const (
	WindowWidth int = 800
	WindowHeight int = 600
)
const Gravity float64 = 0.35
const CoyoteTime time.Duration = time.Millisecond * 7

//type CollisionType int8
//test comment

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
	State          PlayerState//func()
	StateStartTime time.Time
	Mods PlayerModuleList
}

type PlayerState func()

type PhysObject struct {
	Body      *entities.Moving
	ActiColls ActiveCollisions
}

type PlayerModuleList struct {
	WallJump PlayerModule
}

type PlayerModule struct {
	Equipped bool
	Obtained bool
}

func (p *Player) AirState() { //start in air state
	//print("a")

	if player.PhysObject.ActiColls.GroundHit {
		p.SetState(p.GroundState)
		return
	} else {
		if p.Mods.WallJump.Equipped {
			if p.PhysObject.ActiColls.LeftWallHit {
				p.SetState(p.WallSlideLeftState)
			} else if p.PhysObject.ActiColls.RightWallHit {
				p.SetState(p.WallSlideRightState)
			}
		}

		p.DoGravity()
	}

	if oak.IsDown(currentControls.Left)  && p.Body.Delta.X() > -AirMaxSpeed {
		// check to prevent inconsistant top speeds
		//(e.g. if you are half a AirAccel away from AirMaxSpeed)
		if p.Body.Delta.X() - AirAccel > -AirMaxSpeed {
			player.Body.Delta.ShiftX(-AirAccel)
		} else {
			p.Body.Delta.SetX(-AirMaxSpeed)
		}
	} else if oak.IsDown(currentControls.Right) && p.Body.Delta.X() < AirMaxSpeed {
		//second verse, same as the first
		if p.Body.Delta.X() + AirAccel < AirMaxSpeed {
			player.Body.Delta.ShiftX(AirAccel)
		} else {
			p.Body.Delta.SetX(AirMaxSpeed)
		}
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
		if oak.IsDown(currentControls.Jump) {
			p.Jump()
		} 

	} else {
		p.SetState(player.CoyoteState)
	}


	if oak.IsDown(currentControls.Left) {
		player.Body.Delta.SetX(-player.Body.Speed.X())
	} else if oak.IsDown(currentControls.Right) {
		player.Body.Delta.SetX(player.Body.Speed.X())
	} else {
		player.Body.Delta.SetX(0)//player.Body.Delta.X()/2)
	}
	//p.Body.Delta.ShiftY(0.001)
	//howIsHittingLabel(p.Body, Ground)
}

//CoyoteState implements "coyote time" a window of time after
//running off an edge in which you can still jump
func (p *Player) CoyoteState() {
	if p.StateStartTime.Add(CoyoteTime).Before(time.Now()) {
		p.SetState(p.AirState)
	}
	//inherit code from AirState
	//p.AirState()
	if oak.IsDown(currentControls.Jump) {
		p.Jump()
	}

}

func (p *Player) WallSlideLeftState() {
	if isJumpInput() {
		p.Body.Delta.SetY(-WallJumpHeight)
		p.Body.Delta.SetX(WallJumpWidth)
		p.SetState(p.WallJumpLaunchState)
		return
	}
	p.AirState()
}

func (p *Player) WallSlideRightState() {
	if isJumpInput() {
		p.Body.Delta.SetY(-WallJumpHeight)
		p.Body.Delta.SetX(-WallJumpWidth)
		p.SetState(p.WallJumpLaunchState)
		return
	}
	p.AirState()
}

//func WallJumpLaunchState is entered after you walljump,
//temporaraly disabling your controls. This should prevent one sided
//walljumps
func (p *Player) WallJumpLaunchState() {
	if p.TimeFromStateStart() >= WallJumpLaunchDuration {
		p.SetState(p.AirState)
		return
	}
	p.DoGravity()
}

func isJumpInput() bool {
	if k, d := oak.IsHeld(currentControls.Jump); k && (d <= time.Millisecond * 30) {
		return true
	} else {
		return false
	}
}

func (p *Player) Jump() {
	p.Body.Delta.ShiftY(-p.Body.Speed.Y())
	p.Body.ShiftY(p.Body.Delta.Y())
	player.SetState(p.AirState)
}

func (p *Player) SetState(state PlayerState) {
	p.StateStartTime = time.Now()
	p.State = state
}

//TimeFromStateStart gets how long it has been since the last state transition
func (p *Player) TimeFromStateStart() time.Duration {
	return time.Now().Sub(p.StateStartTime)
}

func (o *PhysObject) DoGravity() {

	o.Body.Delta.ShiftY(Gravity)
}

func (object *PhysObject) DoCollision(updater func()) {
	_, oldY := object.Body.GetPos()
	updater()
	object.ActiColls = ActiveCollisions{} //reset the struct to be all false

	object.Body.ShiftX(object.Body.Delta.X())
	if hit := collision.HitLabel(object.Body.Space, Ground); hit != nil {
		if object.Body.Delta.X() > 0 { //Right Wall
			object.ActiColls.RightWallHit = true
			object.Body.SetX(hit.X() - object.Body.W)
		} else if object.Body.Delta.X() < 0 { //Left Wall
			object.ActiColls.LeftWallHit = true
			object.Body.SetX(hit.X() + hit.W())
		}
	}

	object.Body.ShiftY(object.Body.Delta.Y())
	if hit := collision.HitLabel(object.Body.Space, Ground); hit != nil {
		if object.Body.Delta.Y() > 0 { //Ground
			object.ActiColls.GroundHit = true
			object.Body.SetY(hit.Y() - object.Body.H)
		} else if object.Body.Delta.Y() < 0 { //Ceiling
			object.ActiColls.CeilingHit = true
			//TODO: make this work like other collision
			object.Body.SetY(oldY)
		}
		object.Body.Delta.SetY(0)
	}

}
//screens are to be stored as json, problebly compressed in the final game
func loadScreen() {
	file, err := os.Open("level.png")
	if err != nil {
		fmt.Print("error when opening screen file:")
		panic(err)
	}
	reader := bufio.NewReader(file)
	conf,_, err := image.DecodeConfig(reader)
	if err != nil {
		fmt.Print("error when decoding screen file config:")
		panic(err)
	}
	//reader.Reset(file)
	levelImage,_, _ := image.Decode(reader)/*&netpbm.DecodeOptions{
		Target:      netpbm.PNM, //this will allow adding more block types later
		Exact:       false,
		PBMMaxValue:1 ,
	})*/
	if err != nil {
		fmt.Print("error when decoding screen file:")
		panic(err)
	}
	var blockArrLenY int = conf.Height
	var blockArrLenX int = conf.Width 
	const blockSize int = 4
	blockArr := make([][]*entities.Solid, blockArrLenX)
	for j := 0; j < blockArrLenX; j++ {
		blockArr[j] = make([]*entities.Solid, blockArrLenY)
		for i := 0; i < blockArrLenY; i++ {
			if levelImage.At(j,i) == color.Black {
				blockArr[j][i] = entities.NewSolid(
					float64(i*blockSize), float64(j*blockSize),
					float64(blockSize), float64(blockSize),
					render.NewColorBox(10, 10,
						color.RGBA{uint8(j), 0, uint8(i), 255}),
					nil, event.CID(j*i+3))
				render.Draw(blockArr[j][i].R, 8+i*j)
			}
		}
	}
}

func loadScene() {
	go loadScreen()

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
	ground3 := entities.NewSolid(300, 200, 20, 500,
		render.NewColorBox(20, 500, color.RGBA{0, 255, 255, 255}),
		nil, 2)


	ground.UpdateLabel(Ground)
	ground2.UpdateLabel(Ground)
	ground3.UpdateLabel(Ground)

	render.Draw(ground.R)
	render.Draw(ground2.R, 1)
	render.Draw(ground3.R, 1)

	//Give player walljump for now
	player.Mods.WallJump.Equipped = true

}

func main() {
	oak.Add("platformer", func(string, interface{}) {
		loadScene()

		player.Body.Bind(func(id int, nothing interface{}) int {

		
			if oak.IsDown(currentControls.Quit) {
				if oak.IsDown(key.I) {
					fmt.Println(player)
				}
				os.Exit(0)
			}

			player.DoCollision(player.State)

			return 0
		}, event.Enter)
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})

	oak.Init("platformer")
}
