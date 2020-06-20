package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/color"

	//"math"
	"os"
	"time"

	//"compress/flate"
	//"gopkg.in/yaml.v2"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"

	"github.com/lolbinarycat/hookshot-oak/level"
	"github.com/lolbinarycat/hookshot-oak/camera"
)

const (
	Ground collision.Label = iota
	NoWallJump
	Death
	Checkpoint
)

var SolidLabels []collision.Label= []collision.Label{
	Ground,
	NoWallJump,
}

const JumpHeight int = 6
const WallJumpHeight float64 = 6
const WallJumpWidth float64 = 3
const WallJumpLaunchDuration time.Duration = time.Millisecond * 230
const (
	AirAccel    float64 = 0.4
	AirMaxSpeed float64 = 3
)
const ClimbSpeed float64 = 3

//Window constants
const (
	WindowWidth  int = 800
	WindowHeight int = 600
)
const Gravity float64 = 0.35

//CoyoteTime is how long CoyoteState lasts
const CoyoteTime time.Duration = time.Millisecond * 7

//JumpInputTime describes the length of time after the jump button is pressed in which it will count as the player jumping.
//Setting this to be too high may result in multiple jumps to occur for one press of the jump button, while setting it too low may result in jumps being eaten.
const JumpInputTime time.Duration = time.Millisecond * 90

//JumpHeightDecTime is how long JumpHeightDecState lasts
const JumpHeightDecTime time.Duration = time.Millisecond * 200

//JsonScreen is a type to unmarshal the json of
//a file with screen (i.e. one screen worth of level) data into
type JsonScreen struct {
	Rects []JsonRect
}

//type JsonRect defines a struct to
//unmarshal json into
type JsonRect struct {
	X, Y, W, H float64
	Label      collision.Label //warning: label is hardcoded in json file
}

type ActiveCollisions struct {
	GroundHit    bool
	LeftWallHit  bool
	RightWallHit bool
	CeilingHit   bool
	HLabel,VLabel collision.Label //these define the LAST label that was hit (horizontaly and verticaly), as ints cannot be nil
	LastHitV, LastHitH collision.Space
}

type ControlConfig struct {
	Left, Right, Up, Down, Jump, Climb, Quit string
}

var currentControls ControlConfig = ControlConfig{
	Left:  key.LeftArrow,
	Right: key.RightArrow,
	Up:    key.UpArrow,
	Down:  key.DownArrow,
	Jump:  key.Z,
	Climb: key.LeftShift,
	Quit:  key.Q,
}

type Direction uint8

const (
	Left Direction = iota
	Right
)

type Pos struct {
	X float64
	Y float64
}

var player Player

//Player is a type representing the player
//StateInit is a variable that should be set to true when changing states
//it tells the state to initialize values like StateTimer
type Player struct {
	//Body           *entities.Moving
	//ActiColls      ActiveCollisions
	PhysObject
	State          PlayerState //func()
	StateStartTime time.Time
	Mods           PlayerModuleList
	RespawnPos Pos
}

type PlayerState func()

type PhysObject struct {
	Body      *entities.Moving
	ActiColls ActiveCollisions
}

type PlayerModuleList struct {
	WallJump PlayerModule
	Climb    PlayerModule
}

type PlayerModule struct {
	Equipped bool
	Obtained bool
}

//this is the default level for debugLevel, value will be set in loadYamlConfigData()
var debugLevel dlog.Level = dlog.WARN
//var log dlog.Logger = dlog.NewLogger()


func (p *Player) WallJump(dir Direction, EnterLaunch bool) {
	p.Body.Delta.SetY(-WallJumpHeight)

	if dir == Left {
		p.Body.Delta.SetX(-WallJumpWidth)
	} else if dir == Right {
		p.Body.Delta.SetX(WallJumpWidth)
	} else {
		panic("invalid direction to WallJump functon")
	}

	if EnterLaunch {
		p.SetState(p.WallJumpLaunchState)
	} else {
		p.SetState(p.AirState)
	}
}



//DoCliming is the function for shared procceses between
//ClimbRightState and ClimbLeft state
func (p *Player) DoCliming() {
	//this is a hack, and should problebly be fixed
	if int(p.TimeFromStateStart())%2 == 0 && p.ActiColls.LeftWallHit == false && p.ActiColls.RightWallHit == false {
		p.SetState(p.AirState)
	}
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
	if k, d := oak.IsHeld(currentControls.Jump); k && (d <= JumpInputTime) {
		return true
	} else {
		return false
	}
}


func (p *Player) Jump() {
	p.Body.Delta.ShiftY(-p.Body.Speed.Y())
	p.Body.ShiftY(p.Body.Delta.Y())
	player.SetState(p.JumpHeightDecState)
}

func (p *Player) SetState(state PlayerState) {
	p.StateStartTime = time.Now()
	p.State = state
}

func (p *Player) DoAirControls() {
	if oak.IsDown(currentControls.Left) && p.Body.Delta.X() > -AirMaxSpeed {
		// check to prevent inconsistant top speeds
		//(e.g. if you are half a AirAccel away from AirMaxSpeed)
		if p.Body.Delta.X()-AirAccel > -AirMaxSpeed {
			player.Body.Delta.ShiftX(-AirAccel)
		} else {
			p.Body.Delta.SetX(-AirMaxSpeed)
		}
	} else if oak.IsDown(currentControls.Right) && p.Body.Delta.X() < AirMaxSpeed {
		//second verse, same as the first
		if p.Body.Delta.X()+AirAccel < AirMaxSpeed {
			player.Body.Delta.ShiftX(AirAccel)
		} else {
			p.Body.Delta.SetX(AirMaxSpeed)
		}
	}
}

//TimeFromStateStart gets how long it has been since the last state transition
func (p *Player) TimeFromStateStart() time.Duration {
	return time.Now().Sub(p.StateStartTime)
}

func (p *Player) Die() {
	//TODO: death animation
	p.Respawn()
}

func (p *Player) Respawn() {
	p.Body.Delta.SetPos(0,0)
	p.Body.SetPos(player.RespawnPos.X,player.RespawnPos.Y)
}

func (o *PhysObject) DoGravity() {
	o.Body.Delta.ShiftY(Gravity)
}

func (o *PhysObject) DoCustomGravity(grav float64) {
	o.Body.Delta.ShiftY(grav)
}

func (object *PhysObject) DoCollision(updater func()) {
	_, oldY := object.Body.GetPos()
	updater()
	object.ActiColls = ActiveCollisions{} //reset the struct to be all false

	object.Body.ShiftX(object.Body.Delta.X())
	hit := collision.HitLabel(object.Body.Space, SolidLabels...);
	if  hit != nil {
		if object.Body.Delta.X() > 0 { //Right Wall
			object.ActiColls.RightWallHit = true
			object.Body.SetX(hit.X() - object.Body.W)
		} else if object.Body.Delta.X() < 0 { //Left Wall
			object.ActiColls.LeftWallHit = true
			object.Body.SetX(hit.X() + hit.W())
		}
		object.Body.Delta.SetX(0)
		object.ActiColls.HLabel = hit.Label
	}

	object.Body.ShiftY(object.Body.Delta.Y())
	if hit := collision.HitLabel(object.Body.Space, SolidLabels...); hit != nil {
		if object.Body.Delta.Y() > 0 { //Ground
			object.ActiColls.GroundHit = true
			object.Body.SetY(hit.Y() - object.Body.H)
		} else if object.Body.Delta.Y() < 0 { //Ceiling
			object.ActiColls.CeilingHit = true
			//TODO: make this work like other collision
			object.Body.SetY(oldY)
		}
		object.Body.Delta.SetY(0)
		object.ActiColls.VLabel = hit.Label
	}

}
func openFileAsBytes(filename string) ([]byte,error) {
	dlog.Info("opening file",filename)
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	reader := bufio.NewReader(file)
	var byteArr []byte = make([]byte, int(fileSize))
	_, err = reader.Read(byteArr)
	if err != nil {
		return byteArr, err
	}

	return byteArr, nil
}

//TODO: complete this function
func loadYamlConfigData(filename string) {
	dlog.Info("loading yaml config data from",filename)

	rawYaml, err := openFileAsBytes(filename)
	dlog.ErrorCheck(err)
	if err != nil {
		return
	}
	dlog.Verb(rawYaml)

	dlog.Error("function incomplete")

	/*fileInfo, err := file.Stat()
	dlog.ErrorCheck(err)

	if err != nil {
		dlog.Error("unable to get yaml config, using defaults")
		return
	}

	fileSize := fileInfo.Size()
	reader*/
}

//level data is to be stored as json, problebly compressed in the final game
func loadJsonLevelData(filename string) {
	dlog.Info("loading json level data frxom", filename)
	//dlog.Warn("test")
	file, err := os.Open(filename)
	if err != nil {
		//t.Print("error when opening screen file: ")
		//panic(err)
		dlog.Error("error when opening screen file",err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Print("error when getting file info: ")
		panic(err)
	}
	fileSize := fileInfo.Size()
	reader := bufio.NewReader(file)
	var rawJson []byte = make([]byte, int(fileSize))
	_, err = reader.Read(rawJson)
	if err != nil {
		fmt.Print("error when reading file into byte array: ")
		panic(err)
	}
	var screenData JsonScreen
	err = json.Unmarshal(rawJson, &screenData)
	if err != nil {
		defer fmt.Println("json:", rawJson)
		fmt.Print("error unmarshaling screen data: ")
		panic(err)
	}

	for i, rectData := range screenData.Rects {
		rect := entities.NewSolid(rectData.X, rectData.Y, rectData.W, rectData.H,
			render.NewColorBox(int(rectData.W), int(rectData.H), color.RGBA{100, 100, 100, 255}),
			nil, event.CID(i+10))

		rect.UpdateLabel(rectData.Label)
		render.Draw(rect.R)
	}
}

func loadScene() {
	loadJsonLevelData("level.json")

	player.Body = entities.NewMoving(100, 100, 16, 32,
		render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
		nil, 0, 0)
	player.State = player.AirState
	player.RespawnPos = Pos{X : player.Body.X(),Y : player.Body.Y()}

	render.Draw(player.Body.R)
	player.Body.Speed = physics.NewVector(3, float64(JumpHeight))

	level.LoadDevRoom()

	//Give player walljump for now
	player.Mods.WallJump.Equipped = true
	//same for climbing
	player.Mods.Climb.Equipped = true
}



func main() {
	//dlog.SetLogger(log)
	oak.Add("platformer", func(string, interface{}) {
		dlog.SetDebugLevel(debugLevel)
		loadScene()
		oak.ScreenWidth = 800
		oak.ScreenHeight = 600

		camera.StartCameraLoop(player.Body)
		//fmt.Println("screenWidth",oak.ScreenWidth)
		//fmt.Println("screenHeight",oak.ScreenHeight)

		player.Body.Bind(func(id int, nothing interface{}) int {
			//xdlog.SetDebugLevel(dlog.VERBOSE)
			if oak.IsDown(currentControls.Quit) {
				if oak.IsDown(key.I) {
					fmt.Println(player)
				}
				os.Exit(0)
			}
			if player.Body.HitLabel(Checkpoint) != nil {
				player.RespawnPos = Pos{X : player.Body.X(),Y : player.Body.Y()}
			}
			if player.Body.HitLabel(Death) != nil {
				player.Die()
			}

			player.DoCollision(player.State)

			return 0
		}, event.Enter)
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})
	/*err := oak.SetBorderless(true)
	if err != nil {
		panic(err)
	}*/
	//oak.SetViewportBounds(0,0, WindowWidth, WindowHeight)
	//dlog.SetLogLevel()
	oak.Init("platformer")

}
