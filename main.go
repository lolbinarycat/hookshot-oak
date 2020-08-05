package main

import (
	"bufio"
	"image/color"
	"regexp"
	"os"
	"time"

	oak "github.com/oakmound/oak/v2"

	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/physics"
	"github.com/oakmound/oak/v2/render"


	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/lolbinarycat/hookshot-oak/level"
	"github.com/lolbinarycat/hookshot-oak/player"
	"github.com/lolbinarycat/hookshot-oak/physobj"
	_ "github.com/lolbinarycat/hookshot-oak/layers"
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



const HsSpeedX, HsSpeedY = 7, 7
func loadScene() *player.Player {
	plr := loadPlayer()

	plr.Hs.Body.Speed = physics.NewVector(HsSpeedX, HsSpeedY)
	plr.Body.UpdateLabel(labels.Player)
	plr.ExtraSolids = []collision.Label{labels.Block}
	render.Draw(plr.Hs.Body.R, 0)

	// var block PhysObject
	// var block2 PhysObject

	// block.Body = entities.NewMoving(150, 100, 16, 16,
	// 	render.NewColorBox(16, 16, color.RGBA{0, 200, 0, 255}),
	// 	nil, 2, 1)
	// block2.Body = entities.NewMoving(200, 130, 16, 32,
	// 	render.NewColorBox(16, 32, color.RGBA{0, 255, 0, 255}),
	// 	nil, 3, 0)
	// block2.Body.Init()
	// block2.Body.UpdateLabel(labels.Block)
	// render.Draw(block2.Body.R)

	block := physobj.NewBlock(32,32,
		render.NewColorBox(32, 32, color.RGBA{0, 200, 0, 255}),
		300,500+600)
	render.Draw(block.R,1,3)
	
	//block.Body.Init()
	//block.ExtraSolids = []collision.Label{labels.Player}
	//block.Body.UpdateLabel(labels.Block)
	//blocks = append(blocks, &block, &block2)

	//screenSpace = collision.NewSpace(0,0,float64(WindowWidth),float64(WindowHeight),3)

	err := level.LoadTmx("assets/level.tmx")
	if err != nil {
		panic(err)
	}

	player.InitMods(plr)

	//plr.Mods.GiveAll(true)

	dlog.Info("player loaded with data:",*plr)
	return plr
}

const mainSceneName = "platformer"

func main() {
	{
		err := oak.LoadConf("config.json")
		if err != nil {
			dlog.Error("failed to load config.json, error:", err)
		}
		logger := dlog.NewRegexpLogger()
		lvl , err := dlog.ParseDebugLevel(oak.SetupConfig.Debug.Level)
		if err != nil {
			panic(err)
		}
		logger.SetDebugLevel(lvl)
		logger.SetRegexp(regexp.MustCompile(oak.SetupConfig.Debug.Filter))
		dlog.SetLogger(logger)
	}


	// bgLayer = render.NewDynamicHeap()
	// fgLayer = render.NewDynamicHeap()
	// uiLayer = render.NewDynamicHeap()
	// render.SetDrawStack(
	// 	render.NewDynamicHeap(),
	// 	render.NewDynamicHeap(),
	// 	render.NewDynamicHeap(),
	// 	render.NewDynamicHeap(),
	// 	//render.NewDrawFPS(),
	// 	render.NewLogicFPS(),
	// )

	MainSceneStart, MainSceneLoop, MainSceneEnd := buildMainSceneFuncs()
	oak.Add(mainSceneName,
		MainSceneStart,
		MainSceneLoop,
		MainSceneEnd,
	)

	oak.AddScene("titlescreen",titlescreenScene)

	BindCommands()
	dlog.Info("Commands bound")

	oak.SetupConfig.Screen = oak.Screen{Height: 600, Width: 800}
	oak.SetupConfig.FrameRate = 60
	oak.SetAspectRatio(8.0 / 6.0)
	oak.Init("titlescreen")
}

