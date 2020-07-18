package main

import (
	"bufio"
	"image/color"

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
	"github.com/oakmound/oak/v2/physics"
	"github.com/oakmound/oak/v2/render"



	"github.com/lolbinarycat/hookshot-oak/collectables"
	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/lolbinarycat/hookshot-oak/level"
	"github.com/lolbinarycat/hookshot-oak/player"

	"github.com/lolbinarycat/utils"
)

const Frame = time.Second / 60

const RunSpeed float64 = 2.8

//Window constants
const (
	WindowWidth  int = 800
	WindowHeight int = 600
)

var loadSave = false

type Player player.Player
type PhysObject = player.PhysObject
type PlayerModule player.PlayerModule
type PlayerState player.PlayerState

//var block PhysObject //this is global temporaraly

//type Body *entities.Moving

//this is the default level for debugLevel,
//value will be set in loadYamlConfigData()
var debugLevel dlog.Level = /** dlog.VERBOSE /*/ dlog.INFO /**/

//temporary global
var blocks []*PhysObject

//var log dlog.Logger = dlog.NewLogger()

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

//TODO: complete this function (don't)
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

const PlayerWidth = 12
const PlayerHeight = 12

const HsWidth = 4
const HsHeight = 4

func loadPlayer() *player.Player {
	var eyeColor = color.RGBA{0, 255, 255, 255}
	playerSprite := utils.Check2(
		render.LoadSprite("assets/images", "player_new.png")).(render.Renderable)

	var plr = new(player.Player)
	eye1 := render.NewColorBox(1, 4, eyeColor)
	eye2 := eye1.Copy().(*render.Sprite)
	plr.Eyes = [2]*render.Sprite{eye1, eye2}
	plr.Body = entities.NewMoving(300, 400, PlayerWidth, PlayerHeight,
		playerSprite,
		nil, 0, 0)
	plr.Body.Init()
	plr.Body.Space.UpdateLabel(labels.Player)

	eye1.LayeredPoint.Vector = eye1.Attach(plr.Body, 4, 3)
	eye2.LayeredPoint.Vector = eye1.Attach(plr.Body, 8, 3)

	player.SetPlayer(0, plr)

	render.Draw(eye1, 2)
	render.Draw(eye2, 2)

	plr.State = player.RespawnFallState
	plr.RespawnPos = player.Pos{X: plr.Body.X(), Y: plr.Body.Y()}
	render.Draw(plr.Body.R, 1)
	plr.Body.Speed = physics.NewVector(3, float64(player.JumpHeight))

	plr.Hs.Body = &*entities.NewMoving(100, 100, HsWidth, HsHeight,
		render.NewColorBox(HsHeight, HsWidth, color.RGBA{0, 0, 255, 255}),
		nil, 1, 0)
	plr.Hs.Body.Init()

	return plr
}

func loadScene() *player.Player {
	//loadJsonLevelData("level.json")

	//AttachMut(&eye1.LayeredPoint.Vector,plr.Body)
	//vectAttach(eye1).AttachMut(plr.Body)
	plr := loadPlayer()

	plr.Hs.Body.Speed = physics.NewVector(3, 3)
	plr.Body.UpdateLabel(labels.Player)
	plr.ExtraSolids = []collision.Label{labels.Block}
	//plr.Hs.Body = entities.NewInteractive(100, 10, 4, 4,
	//	render.NewColorBox(16, 16, color.RGBA{0, 0, 255, 255}),
	//	nil, 1, 0)

	//plr.Body.Doodad.Point.Attach(plr.Hs.Body)
	//plr.Body.AttachX(plr.Hs.Body,0)
	render.Draw(plr.Hs.Body.R, 0)

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

	err := level.LoadDevRoom()
	if err != nil {
		panic(err)
	}

	player.InitMods(plr)
	if loadSave {
		err := plr.Load("save.json")
		if err != nil {
			panic(err)
		}
	}

	plr.Mods.GiveAll(true)
	modClct := collectables.NewModuleClct(120, 550, 8, 8,
		render.NewColorBox(8, 8, color.RGBA{0, 255, 100, 255}), 72, "hs")
	render.Draw(modClct.React.R, 3)

	//render.NewDrawFPS()
	//render.Draw(fps)
	return plr
}

//var progStartTime time.Time
func main() {
	// Apperenly 1 DynamicHeap = 1 layer.
	render.SetDrawStack(
		render.NewDynamicHeap(),
		render.NewDynamicHeap(),
		render.NewDynamicHeap(),
		render.NewDynamicHeap(),
		//render.NewDrawFPS(),
		render.NewLogicFPS(),
	)
	dlog.SetDebugLevel(debugLevel)

	//MainSceneLoop = func() bool {return true}
	//progStartTime = time.Now()
	//dlog.SetLogger(log)
	//loadPlayer()
	MainSceneStart, MainSceneLoop, MainSceneEnd := buildMainSceneFuncs()
	oak.Add("platformer",
		MainSceneStart,
		MainSceneLoop,
		MainSceneEnd,
	)

	BindCommands()

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
		dlog.Error("failed to load config.json, error:", err)
	}
	oak.SetupConfig.Screen = oak.Screen{Height: 600, Width: 800}
	oak.SetupConfig.FrameRate = 60
	oak.SetAspectRatio(8.0 / 6.0)
	oak.Init("platformer")
}

// defines a playerstate with only a loop function
/*func (p *Player) NewJustLoopState(loopFunc PlayerStateFunc) PlayerState {
	PlayerState{
		Loop:loopFunc,

	}
}*/
