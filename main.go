package main

import (
	"bufio"
	"image/color"

	"os"
	"time"

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



//temporary global
var blocks []*PhysObject



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
	plr := loadPlayer()

	plr.Hs.Body.Speed = physics.NewVector(3, 3)
	plr.Body.UpdateLabel(labels.Player)
	plr.ExtraSolids = []collision.Label{labels.Block}
	render.Draw(plr.Hs.Body.R, 0)

	var block PhysObject
	var block2 PhysObject

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

	err := level.LoadTmx("assets/level.tmx")
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

	return plr
}

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

	MainSceneStart, MainSceneLoop, MainSceneEnd := buildMainSceneFuncs()
	oak.Add("platformer",
		MainSceneStart,
		MainSceneLoop,
		MainSceneEnd,
	)

	BindCommands()

	err := oak.LoadConf("config.json")
	if err != nil {
		dlog.Error("failed to load config.json, error:", err)
	}
	oak.SetupConfig.Screen = oak.Screen{Height: 600, Width: 800}
	oak.SetupConfig.FrameRate = 60
	oak.SetAspectRatio(8.0 / 6.0)
	oak.Init("platformer")
}

