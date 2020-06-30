package main

import (
	"bufio"
	"fmt"
	"image/color"
	"math"
	"strconv"

	//"math"
	"os"
	"time"

	//"compress/flate"
	//"gopkg.in/yaml.v2"

	//"github.com/disintegration/gift"
	oak "github.com/oakmound/oak/v2"

	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/physics"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"

	"github.com/lolbinarycat/hookshot-oak/camera"
	"github.com/lolbinarycat/hookshot-oak/direction"
	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/lolbinarycat/hookshot-oak/level"

	"github.com/lolbinarycat/utils"
)

const JumpHeight int = 5
const WallJumpHeight float64 = 6
const WallJumpWidth float64 = 3

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



//JumpInputTime describes the length of time after the jump button is pressed in which it will count as the player jumping.
//Setting this to be too high may result in multiple jumps to occur for one press of the jump button, while setting it too low may result in jumps being eaten.
const JumpInputTime time.Duration = time.Millisecond * 90

const HsInputTime time.Duration = time.Millisecond * 70


type ActiveCollisions struct {
	GroundHit          bool
	LeftWallHit        bool
	RightWallHit       bool
	CeilingHit         bool
	HLabel, VLabel     collision.Label // these define the LAST label that was hit (horizontaly and verticaly), as ints cannot be nil
	LastHitV, LastHitH event.CID // cid of the last collision space that was hit
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

var curCtrls = &currentControls

//type Direction uint8

// const (
//	 Left Direction = iota
//	 Right
// )

type Pos struct {
	X float64
	Y float64
}

//var block PhysObject //this is global temporaraly

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
	HeldObj        *entities.Moving
	Eyes  [2]*render.Sprite
}

type Hookshot struct {
	PhysObject
	X, Y   float64
	Active bool
}

type PlayerState struct {
	Start, Loop, End PlayerStateFunc
}

type PlayerStateFunc func(*Player)

type PhysObject struct {
	Body      *entities.Moving
	ActiColls ActiveCollisions
	//ExtraSolids defines labels that should be solid only for this object.
	ExtraSolids []collision.Label
}

//type Body *entities.Moving

type PlayerModuleList struct {
	WallJump  PlayerModule
	Climb     PlayerModule
	Hookshot  PlayerModule
	BlockPush PlayerModule
	BlockPull,
	Fly,
	GroundPound, // FloorDollar
	GroundPoundJump,
	HsItemGrab PlayerModule
}

type PlayerModule struct {
	Equipped bool
	Obtained bool
}

//whether modules should be automaticaly equipped when recived
var autoEquipMods bool = true

//this is the default level for debugLevel,
//value will be set in loadYamlConfigData()
var debugLevel dlog.Level = /** dlog.VERBOSE /*/ dlog.INFO/**/ 

//temporary global
var blocks []*PhysObject

//var log dlog.Logger = dlog.NewLogger()

func (p *Player) WallJump(dir direction.Dir, EnterLaunch bool) {
	p.Body.Delta.SetY(-WallJumpHeight)

	if dir.IsLeft() {
		p.Body.Delta.SetX(-WallJumpWidth)
	} else if dir.IsRight() {
		p.Body.Delta.SetX(WallJumpWidth)
	} else {
		panic("invalid direction to WallJump functon")
	}

	if EnterLaunch {
		p.SetState(WallJumpLaunchState)
	} else {
		p.SetState(AirState)
	}
}

//DoCliming is the function for shared procceses between
//ClimbRightState and ClimbLeft state
func (p *Player) DoCliming() {
	//this is a hack, and should problebly be fixed
	if int(p.TimeFromStateStart())%2 == 0 && p.ActiColls.LeftWallHit == false && p.ActiColls.RightWallHit == false {
		p.SetState(AirState)
	}
	if !oak.IsDown(currentControls.Climb) {
		p.SetState(AirState)
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
		p.SetState(HsStartState)
	}
}

func (p *Player) Jump() {
	p.Body.Delta.ShiftY(-p.Body.Speed.Y())
	p.Body.ShiftY(p.Body.Delta.Y())
	player.SetState(JumpHeightDecState)
}

func (p *Player) SetState(state PlayerState) {
	defer func() {
		if r := recover(); r != nil {
			dlog.Error("error while setting state", r)
			p.State = state
		}
	}()

	p.State.End(p)
	p.StateStartTime = time.Now()

	p.State = state
	p.State.Start(p)
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
	p.SetState(RespawnFallState)
	p.Body.Delta.SetPos(0, 0)
	p.Body.SetPos(player.RespawnPos.X, player.RespawnPos.Y)
}

func (o *PhysObject) DoGravity() {
	o.Body.Delta.ShiftY(Gravity)
}

func (o *PhysObject) DoCustomGravity(grav float64) {
	o.Body.Delta.ShiftY(grav)
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

var screenSpace *collision.Space

func loadScene() {
	//loadJsonLevelData("level.json")
	var eyeColor = color.RGBA{0,255,255,255}
	playerSprite := utils.Check2(
		render.LoadSprite("assets/images","player_new.png")).(render.Renderable)

	eye1 := render.NewColorBox(1,4,eyeColor)
	eye2 := eye1.Copy().(*render.Sprite)
	player.Eyes = [2]*render.Sprite{eye1,eye2}
	player.Body = entities.NewMoving(100, 100, 12, 12,
		playerSprite,
		nil, 0, 0)
	player.Body.Init()
	eye1.LayeredPoint.Vector = eye1.Attach(player.Body,4,3)
	eye2.LayeredPoint.Vector = eye1.Attach(player.Body,8,3)
	//AttachMut(&eye1.LayeredPoint.Vector,player.Body)
	//vectAttach(eye1).AttachMut(player.Body)
	render.Draw(eye1,1)
	render.Draw(eye2,1)

	player.State = RespawnFallState
	player.RespawnPos = Pos{X: player.Body.X(), Y: player.Body.Y()}
	render.Draw(player.Body.R)
	player.Body.Speed = physics.NewVector(3, float64(JumpHeight))

	player.Hs.Body = entities.NewMoving(100, 100, 4, 4,
		render.NewColorBox(4, 4, color.RGBA{0, 0, 255, 255}),
		nil, 1, 0)
	player.Hs.Body.Init()

	player.Hs.Body.Speed = physics.NewVector(3, 3)
	player.Body.UpdateLabel(labels.Player)
	player.ExtraSolids = []collision.Label{labels.Block}
	//player.Hs.Body = entities.NewInteractive(100, 10, 4, 4,
	//	render.NewColorBox(16, 16, color.RGBA{0, 0, 255, 255}),
	//	nil, 1, 0)

	//player.Body.Doodad.Point.Attach(player.Hs.Body)
	//player.Body.AttachX(player.Hs.Body,0)
	render.Draw(player.Hs.Body.R,-1)

	var block PhysObject
	var block2 PhysObject
	//(mod.ResizeToFit(16*2,16*2, gift.NearestNeighborResampling))//.(*render.Sprite)
	//if err != nil {
	//	panic(err)
	//}
	block.Body = entities.NewMoving(150, 100, 16, 16,
		render.NewColorBox(16, 16, color.RGBA{0, 200, 0, 255}),
		nil, 2, 1)
	block2.Body = entities.NewMoving(200, 130, 16, 32,
		render.NewColorBox(16, 32, color.RGBA{0, 255, 0, 255}),
		nil, 3, 0)
	block2.Body.Init()
	block2.Body.UpdateLabel(labels.Block)
	render.Draw(block2.Body.R)

	render.Draw(block.Body.R)
	block.Body.Init()
	block.ExtraSolids = []collision.Label{labels.Player}
	block.Body.UpdateLabel(labels.Block)
	blocks = append(blocks, &block, &block2)


	//screenSpace = collision.NewSpace(0,0,float64(WindowWidth),float64(WindowHeight),3)

	level.LoadDevRoom()

	//Give player walljump for now
	//player.Mods.WallJump.Equipped = true
	//same for climbing
	//player.Mods.Climb.Equipped = true
	// " "
	//player.Mods.Hookshot.Equipped = true
	//player.Mods.BlockPush.Equipped = true
	{
		m := &player.Mods
		GiveMods(&m.BlockPush,
			&m.Climb,
			&m.Hookshot,
			&m.WallJump,
			&m.BlockPull,
			&m.HsItemGrab,
			&m.GroundPound,
			&m.GroundPoundJump,
			&m.Fly)
	}
	render.NewDrawFPS()
	//render.Draw(fps)
}

//var progStartTime time.Time
func main() {
	dlog.SetDebugLevel(debugLevel)
	initStates()
	//progStartTime = time.Now()
	//dlog.SetLogger(log)
	oak.Add("platformer", func(string, interface{}) {
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
				//oak.MoveWindow(20, 20, 800, 600)
				//oak.SetAspectRatio(16 / 9)

			}
			if oak.IsDown(currentControls.Quit) {
				if oak.IsDown(key.I) {
					fmt.Println(player)
				}
				os.Exit(0)
			}

			if player.Body.HitLabel(labels.Checkpoint) != nil {
				player.RespawnPos = Pos{X: player.Body.X(), Y: player.Body.Y()}
			}
			if player.Body.HitLabel(labels.Death) != nil {
				player.Die()
			}

			//blocks := collision.WithLabels(labels.Block)
			for _, block := range blocks {
				block.DoCollision(block.BlockUpdater)
				//	block.CID.E().(PhysObject).DoCollision(block.BlockUpdater)
			}

			player.DoCollision(player.DoStateLoop)

			if !player.Hs.Active {
				player.Hs.Body.SetPos(player.Body.X()+hsOffX, //+player.Hs.X,
					player.Body.Y()+hsOffY) //+player.Hs.Y)
			}

			player.Hs.DoCollision(HsUpdater)

			//player.Eyes[1].SetX(5)

			return 0
		}, event.Enter)
	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})

	oak.AddCommand("kill", func(args []string) {
		player.Die()
	})
	oak.AddCommand("setAspectRatio", func(args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: setAspectRatio .75")
		}
		xToY, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			fmt.Println("failed to parse floating point ratio")
		}
		oak.SetAspectRatio(xToY)
		oak.ChangeWindow(oak.ScreenWidth, oak.ScreenHeight)
	})
	oak.AddCommand("fly", func(args []string) {
		if player.Mods.Fly.Equipped {
			if len(args) == 1 &&
				utils.EqualsAny(args[0],"stop","end","halt") {

				player.SetState(AirState)
			} else {
				player.SetState(FlyState)
			}
		}
	})


	/*err := oak.SetBorderless(true
	/*err := oak.SetBorderless(true)
	if err != nil {
		panic(err)
	}*/
	//dlog.SetLogLevel()
	//oak.SetAspectRatio(float64(6/8))
	//oak.ScreenWidth = 800
	//oak.ScreenHeight = 600
	err := oak.LoadConf("config.json")
	if err != nil {
		dlog.Error("failed to load config.json, error:",err)
	}
	oak.SetupConfig.Screen = oak.Screen{Height:600,Width:800}
	oak.SetAspectRatio(8.0/6.0)
	oak.Init("platformer")
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
	p.Hs.Body.Delta.SetPos(0, 0)
	p.SetState(AirState)
}

func GiveMods(mods ...*PlayerModule) {
	for _, m := range mods {
		m.Obtained = true
		if autoEquipMods {
			m.Equipped = true
		}
	}
}

func (b *PhysObject) BlockUpdater() {
	//b.Body.ApplyFriction(1)
	//b.Body.Delta.
	b.DoGravity()
}

// Depreciated
func (p *Player) GrabObject(xOff, yOff, maxDist float64, targetLabels ...collision.Label) (bool, event.CID) {
	if len(targetLabels) > 1 {
		dlog.Error("muliple labels not implemented yet")
	}

	id, ent := event.ScanForEntity(func(e interface{}) bool {
		if ent, ok := e.(*entities.Moving); ok {

			if ent.Space.Label != targetLabels[0] {
				dlog.Verb("label check failed")
				return false
			}
			if !(ent.Space.CID == p.ActiColls.LastHitH) {
				dlog.Verb("id is equal. id:", ent.CID)
				return false
			}

			if ent.DistanceTo(p.Body.X()+xOff, p.Body.Y()+yOff) <=
				maxDist+(math.Max(ent.W, ent.H)) {

				dlog.Verb("distance condition fufilled")
				// if the entity has the correct label, and is within the max distance:
				return true
			}

			//dlog.Verb("d ==",d)
		} else {
			// if the entity is not a entities.Solid, we cannot grab it
			dlog.Verb("type check failed")
			return false
		}
		//this is just to stop "missing return at end of function"
		return false
	})

	// if id is equal to -1, it means ScanForEntity was unable
	// to find an entity within the given paramaters
	if id == -1 {
		dlog.Verb("ScanForEntity Failed")
		return false, -1
	}
	//p.HeldObjId = event.CID(id)
	if mov, ok := ent.(*entities.Moving); ok {
		p.HeldObj = &*mov
		dlog.Verb("HeldObj set")
	} else {
		dlog.Verb("ent exists, but is not *entities.Moving")
		return false, -1
	}

	return true, event.CID(id)
}

func (p *Player) GrabObjRight(targetLabels ...collision.Label) (bool, event.CID) {
	return p.GrabObject(p.Body.W, p.Body.H, p.Body.W, targetLabels...)
}

func (p *Player) GrabObjLeft(targetLabels ...collision.Label) (bool, event.CID) {
	return p.GrabObject(-p.Body.W, -p.Body.H, p.Body.W, targetLabels...)
}

// GetLastHitObj attempts to get an entity from a PhysObject's
// ActiColls.LastHit* attribute. .LastHitH if Horis == true,
// and .LastHitV if false.
// it will return nil if unsucssesful.

func (o *PhysObject) GetLastHitObj(Horis bool) *entities.Moving {
	_, iface := event.ScanForEntity(func(ent interface{}) bool {
		mov, ok := ent.(*entities.Moving)
		if !ok {
			return false
		}
		if (Horis && mov.Space.CID == o.ActiColls.LastHitH) ||
			(!Horis && mov.Space.CID == o.ActiColls.LastHitV) {
			return true
		}
		return false
	})
	return iface.(*entities.Moving)
}



// defines a playerstate with only a loop function
/*func (p *Player) NewJustLoopState(loopFunc PlayerStateFunc) PlayerState {
	PlayerState{
		Loop:loopFunc,

	}
}*/

func (p *Player) DoStateLoop() {
	p.State.Loop(p)
}

func (o *PhysObject) IsWallHit() bool {
	if o.ActiColls.LeftWallHit || o.ActiColls.RightWallHit {
		return true
	}
	return false
}

// Depreciated
func (p *Player) IsHsInPlayer() bool {
	xover, yover := p.Hs.Body.Space.Overlap(p.Body.Space)
	if xover >= p.Hs.Body.W || yover >= p.Hs.Body.H {
		return true
	}
	return false
}

func (p *Player) DoHsCheck() bool {
	if p.IsHsInPlayer() || p.IsWallHit() || p.Hs.X <= 0 {
		p.EndHs()
		return true
	}
	return false
}

//func (p *Player) SetEyePos() {
//
//}



//func AttachMut(a *vectAttach,attachTo physics.Attachable,offsets... float64) {
//	*a = a.Attach(attachTo,offsets...)
//}


func (o *PhysObject) HasHitInDir(dir direction.Dir) bool {
	return (dir.IsRight() && o.ActiColls.RightWallHit) ||
		(dir.IsLeft() && o.ActiColls.LeftWallHit)   ||
		(dir.IsUp() && o.ActiColls.CeilingHit)      ||
		(dir.IsDown() && o.ActiColls.GroundHit)
}
