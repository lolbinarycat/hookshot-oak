package main

import (
	"bufio"
	"fmt"
	"image/color"

	//"math"
	"os"
	"time"

	//"compress/flate"
	//"gopkg.in/yaml.v2"

	"github.com/oakmound/oak/v2"

	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/physics"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"

	"github.com/lolbinarycat/hookshot-oak/camera"
	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/lolbinarycat/hookshot-oak/level"
)

const JumpHeight int = 6
const WallJumpHeight float64 = 6
const WallJumpWidth float64 = 3
const WallJumpLaunchDuration time.Duration = time.Millisecond * 230
const (
	AirAccel    float64 = 0.4
	AirMaxSpeed float64 = 3
)
const ClimbSpeed float64 = 3
const RunSpeed float64 = 2.8

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

const HsInputTime time.Duration = time.Millisecond * 70

//JumpHeightDecTime is how long JumpHeightDecState lasts
const JumpHeightDecTime time.Duration = time.Millisecond * 200

const HsExtendTime time.Duration = time.Second * 2

type ActiveCollisions struct {
	GroundHit          bool
	LeftWallHit        bool
	RightWallHit       bool
	CeilingHit         bool
	HLabel, VLabel     collision.Label //these define the LAST label that was hit (horizontaly and verticaly), as ints cannot be nil
	LastHitV, LastHitH collision.Space
}

type ControlConfig struct {
	Left, Right, Up, Down, Jump, Hs, Climb, Quit string
}

var currentControls ControlConfig = ControlConfig{
	Left:  key.LeftArrow,
	Right: key.RightArrow,
	Up:    key.UpArrow,
	Down:  key.DownArrow,
	Jump:  key.Z,
	Hs:    key.X,
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

var block PhysObject //this is global temporaraly

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
	RespawnPos     Pos
	Hs             Hookshot
}

type Hookshot struct {
	PhysObject
	X, Y   float64
	Active bool
}

type PlayerState func()

type PhysObject struct {
	Body      *entities.Moving
	ActiColls ActiveCollisions
	//ExtraSolids defines labels that should be solid only for this object.
	ExtraSolids []collision.Label
}

type PlayerModuleList struct {
	WallJump PlayerModule
	Climb    PlayerModule
	Hookshot PlayerModule
}

type PlayerModule struct {
	Equipped bool
	Obtained bool
}

//this is the default level for debugLevel,
//value will be set in loadYamlConfigData()
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

	p.StateCommon()
}

func isJumpInput() bool {
	return isButtonPressedWithin(currentControls.Jump, JumpInputTime)
}

func isButtonPressedWithin(button string, dur time.Duration) bool {
	if k, d := oak.IsHeld(button); k && (d <= dur) {
		return true
	} else {
		return false
	}
}

func isHsInput() bool {
	return isButtonPressedWithin(currentControls.Hs, HsInputTime)
}

func (p *Player) ifHsPressedStartHs() {
	if isHsInput() {
		p.SetState(p.HsStartState)
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
	p.SetState(p.RespawnFallState)
	p.Body.Delta.SetPos(0, 0)
	p.Body.SetPos(player.RespawnPos.X, player.RespawnPos.Y)
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
	hit := collision.HitLabel(object.Body.Space,
		append(labels.Solids, object.ExtraSolids...)...)
	if hit != nil {
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
	if hit := collision.HitLabel(object.Body.Space,
		append(object.ExtraSolids, labels.Solids...)...); hit != nil {
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
func openFileAsBytes(filename string) ([]byte, error) {
	dlog.Info("opening file", filename)
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
	dlog.Info("loading yaml config data from", filename)

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

func loadScene() {
	//loadJsonLevelData("level.json")

	player.Body = entities.NewMoving(100, 100, 16, 16,
		render.NewColorBox(16, 16, color.RGBA{255, 0, 0, 255}),
		nil, 0, 0)
	player.State = player.RespawnFallState
	player.RespawnPos = Pos{X: player.Body.X(), Y: player.Body.Y()}
	render.Draw(player.Body.R)
	player.Body.Speed = physics.NewVector(3, float64(JumpHeight))

	player.Hs.Body = entities.NewMoving(100, 100, 4, 4,
		render.NewColorBox(4, 4, color.RGBA{0, 0, 255, 255}),
		nil, 1, 0)

	player.Hs.Body.Speed = physics.NewVector(3, 3)
	player.Body.UpdateLabel(labels.Player)
	player.ExtraSolids = []collision.Label{}
	//player.Hs.Body = entities.NewInteractive(100, 10, 4, 4,
	//	render.NewColorBox(16, 16, color.RGBA{0, 0, 255, 255}),
	//	nil, 1, 0)

	//player.Body.Doodad.Point.Attach(player.Hs.Body)
	//player.Body.AttachX(player.Hs.Body,0)
	render.Draw(player.Hs.Body.R)

	level.LoadDevRoom()

	block.Body = entities.NewMoving(150, 100, 16, 16,
		render.NewColorBox(16, 16, color.RGBA{0, 200, 0, 255}),
		nil, 780, 0)
	render.Draw(block.Body.R)
	block.ExtraSolids = []collision.Label{labels.Player}
	block.Body.UpdateLabel(labels.Block)

	//Give player walljump for now
	player.Mods.WallJump.Equipped = true
	//same for climbing
	player.Mods.Climb.Equipped = true
	// " "
	player.Mods.Hookshot.Equipped = true
}

func main() {
	//dlog.SetLogger(log)
	oak.Add("platformer", func(string, interface{}) {
		dlog.SetDebugLevel(debugLevel)
		loadScene()

		camera.StartCameraLoop(player.Body)
		//fmt.Println("screenWidth",oak.ScreenWidth)
		//fmt.Println("screenHeight",oak.ScreenHeight)

		hsOffX := player.Body.W/2 - player.Hs.Body.H/2
		hsOffY := player.Body.H/2 - player.Hs.Body.H/2

		player.Body.Bind(func(id int, nothing interface{}) int {
			//xdlog.SetDebugLevel(dlog.VERBOSE)
			if oak.IsDown(key.L) {
				//oak.ScreenWidth = 800
				//oak.ScreenHeight = 600
				//oak.ChangeWindow(800,600)
				oak.MoveWindow(20, 20, 800, 600)
				oak.SetAspectRatio(16 / 9)

			}
			if oak.IsDown(currentControls.Quit) {
				if oak.IsDown(key.I) {
					fmt.Println(player)
					fmt.Println("block:\n", block)
				}
				os.Exit(0)
			}
			if player.Body.HitLabel(labels.Checkpoint) != nil {
				player.RespawnPos = Pos{X: player.Body.X(), Y: player.Body.Y()}
			}
			if player.Body.HitLabel(labels.Death) != nil {
				player.Die()
			}

			block.DoCollision(block.BlockUpdater)

			player.DoCollision(player.State)

			if !player.Hs.Active {
				player.Hs.Body.SetPos(player.Body.X()+hsOffX, //+player.Hs.X,
					player.Body.Y()+hsOffY) //+player.Hs.Y)
			}

			player.Hs.DoCollision(HsUpdater)

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
	//dlog.SetLogLevel()
	oak.Init("platformer")
	oak.UseAspectRatio = true

}

func HsUpdater() {
	hsOffX := player.Body.W/2 - player.Hs.Body.H/2
	hsOffY := player.Body.H/2 - player.Hs.Body.H/2

	//set hookshot's relitive position to be accurate
	player.Hs.X = player.Hs.Body.X() - player.Body.X() - hsOffX
	player.Hs.Y = player.Hs.Body.Y() - player.Body.Y() - hsOffY
}

func (p *Player) EndHs() {
	p.Hs.Active = false
	p.Hs.X = 0
	p.Hs.Y = 0
	p.SetState(p.AirState)
}

func (b *PhysObject) BlockUpdater() {
	hit := collision.HitLabel(b.Body.Space, labels.Player)

	if hit != nil {
		xover, _ := b.Body.Space.Overlap(hit)
		// xover < 8 {
		b.Body.Delta.SetX(xover)
		//

	} 

	/*if (b.ActiColls.LeftWallHit || b.ActiColls.RightWallHit) ||
		(b.ActiColls.HLabel == labels.Player) {
		b.Body.Delta.SetX(player.Body.Delta.X())
	}*/
	if b.ActiColls.GroundHit == false {
		b.DoGravity()
	}
}
