package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/color"
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
const ClimbSpeed float64 = 3
//Window constants
const (
	WindowWidth int = 800
	WindowHeight int = 600
)
const Gravity float64 = 0.35
const CoyoteTime time.Duration = time.Millisecond * 7

//JsonScreen is a type to unmarshal the json of
//a file with screen (i.e. one screen worth of level) data into
type JsonScreen struct {
	Rects []JsonRect
}
//type JsonRect defines a struct to
//unmarshal json into
type JsonRect struct {
	X,Y,W,H float64
	Label collision.Label //warning: label is hardcoded in json file
}

type ActiveCollisions struct {
	GroundHit    bool
	LeftWallHit  bool
	RightWallHit bool
	CeilingHit   bool
}

type ControlConfig struct {
	Left, Right, Up, Down, Jump, Climb, Quit string
}

var currentControls ControlConfig = ControlConfig{
	Left:  key.LeftArrow,
	Right: key.RightArrow,
	Up: key.UpArrow,
	Down: key.DownArrow,
	Jump:  key.Z,
	Climb: key.LeftShift,
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
	Climb PlayerModule
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
		if isJumpInput() {
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
	if isJumpInput() {
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
	if oak.IsDown(currentControls.Climb) && p.Mods.Climb.Equipped {
		p.SetState(p.ClimbLeftState)
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
	if oak.IsDown(currentControls.Climb) && p.Mods.Climb.Equipped {
		p.SetState(p.ClimbRightState)
		return //return to stop airstate for overwriting our change
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

func (p *Player) ClimbRightState() {
	if isJumpInput() {
		p.SetState(p.WallSlideRightState)
		p.WallSlideRightState()
	}
	p.DoCliming()
}

func (p *Player) ClimbLeftState() {
	if isJumpInput() {
		p.SetState(p.WallSlideLeftState)
		p.WallSlideRightState()
	}
	p.DoCliming()
}

//DoCliming is the function for shared procceses between
//ClimbRightState and ClimbLeft state
func (p *Player) DoCliming() {
	if !oak.IsDown(currentControls.Climb) {
		p.SetState(p.AirState)
	}
	if oak.IsDown(currentControls.Up) {
		p.Body.Delta.SetY(-ClimbSpeed)
	} else if oak.IsDown(currentControls.Down) {
		p.Body.Delta.SetY(ClimbSpeed)
	} else {
		p.Body.Delta.SetY(0)
	}
}

func isJumpInput() bool {
	if k, d := oak.IsHeld(currentControls.Jump); k && (d <= time.Millisecond * 100) {
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
func loadJsonScreen(filename string) {
	file, err := os.Open("level.json")
	if err != nil {
		fmt.Print("error when opening screen file: ")
		panic(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Print("error when getting file info: ")
		panic(err)
	}
	fileSize := fileInfo.Size()
	reader := bufio.NewReader(file)
	var rawJson []byte = make([]byte,int(fileSize))
	_,err = reader.Read(rawJson)
	if err != nil {
		fmt.Print("error when reading file into byte array: ")
		panic(err)
	}
	var screenData JsonScreen
	err = json.Unmarshal(rawJson, &screenData)
	if err != nil {
		defer fmt.Println("json:",rawJson)
		fmt.Print("error unmarshaling screen data: ")
		panic(err)
	}

	for i , rectData := range screenData.Rects {
		rect := entities.NewSolid(rectData.X, rectData.Y, rectData.W, rectData.H,
			render.NewColorBox(int(rectData.W), int(rectData.H), color.RGBA{100, 100, 100, 255}),
			nil, event.CID(i+10) )
		
		rect.UpdateLabel(rectData.Label)
		render.Draw(rect.R)
	}
}

func loadScene() {
	loadJsonScreen("level.json")

	player.Body = entities.NewMoving(100, 100, 16, 32,
		render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
		nil, 0, 0)
	player.State = player.AirState

	render.Draw(player.Body.R)
	player.Body.Speed = physics.NewVector(3, float64(JumpHeight))

	ground := entities.NewSolid(10, 400, 500, 20,
		render.NewColorBox(500, 20, color.RGBA{0, 0, 255, 255}),
		nil, 0)
	ground2 := entities.NewSolid(40, 200, 20, 500,
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
	//same for climbing
	player.Mods.Climb.Equipped = true
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
			//oak.SetScreen(0,0)
			player.DoCollision(player.State)

			return 0
		}, event.Enter)
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})
	oak.SetAspectRatio(8/6)
	oak.SetFullScreen(true)
	oak.SetViewportBounds(0,0, 800, 600)
	oak.Init("platformer")

	oak.ChangeWindow(800,600)

}
