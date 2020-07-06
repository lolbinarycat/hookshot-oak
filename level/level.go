package level

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/color"
	"os"

	//"strconv"

	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"

	//"github.com/rustyoz/svg"
	"github.com/lafriks/go-tiled"
	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/lolbinarycat/hookshot-oak/collectables"
)

/*type Rect struct {
	
	//X, 
	//Y,
	W string `xml:"width,attr"`
	H string `xml:"height,attr"`
}*/




//color constants
var (
	Gray    color.RGBA = color.RGBA{100, 100, 100, 255}
	DullRed color.RGBA = color.RGBA{100, 10, 10, 255}
)

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

func newCounter(n int) func() event.CID {
	var num = n - 1
	return func() event.CID {
		num++
		return event.CID(num)
	}
}

func OpenFileAsBytes(filename string) ([]byte, error) {
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

func LoadDevRoom() {
	n := newCounter(3)
	fmt.Println(n())
	// ground := entities.NewSolid(10, 400, 500, 20,
	// 	render.NewColorBox(500, 20, Gray),
	// 	nil, n())
	// ground.Init()

	//wallSprite, err := render.LoadSprite("assets/images","wall.png")

	//dlog.ErrorCheck(err)
	//wallSprite.Modify(mod.Scale(2,2))

	// wall1 := entities.NewSolid(40, 200, 20, 500,
	// 	render.NewColorBox(20, 500, Gray),
	// 	nil, n())
	// wall1.Init()
	// wall2 := entities.NewSolid(300, 200, 20, 500,
	// 	render.NewColorBox(20, 500, Gray),
	// 	nil, n())
	// wall2.Init()
	checkpoint := entities.NewSolid(200, 350, 10, 10,
		render.NewColorBox(10, 10, color.RGBA{0, 0, 255, 255}),
		nil, n())
	checkpoint.Init()
	death := entities.NewSolid(340, 240, 50, 50,
		render.NewColorBox(50, 50, DullRed),
		nil, n())
	death.Init()


	//ground.UpdateLabel(labels.Ground)
	//wall1.UpdateLabel(labels.Ground)
	//wall2.UpdateLabel(labels.Ground)
	checkpoint.UpdateLabel(labels.Checkpoint)
	death.UpdateLabel(labels.Death)

	//render.Draw(ground.R)
	//render.Draw(wall1.R, 1)
	//render.Draw(wall2.R, 1)
	render.Draw(checkpoint.R, 1)
	render.Draw(death.R)

	//LoadJsonLevelData("level.json",-800,0)

	err := LoadTmx("assets/level.tmx")
	if err != nil {
		panic(err)
	}

	//err := loadSvg("level.svg", -800, 0)
	//dlog.ErrorCheck(err)
}

// func loadSvg(filename string, offsetX, offsetY float64) error {

// 	dlog.Info("opening file", filename)
// 	file, err := os.Open(filename)
// 	defer file.Close()
// 	if err != nil {
// 		return err
// 	}
// 	/*fileInfo, err := file.Stat()
// 	if err != nil {
// 		return nil, err
// 	}
// 	fileSize := fileInfo.Size()*/

// 	/*

// 	svgData, err := svg.ParseSvgFromReader(reader, filename, 1)
// 	*/
// 	reader := bufio.NewReader(file)
// 	decoder := xml.NewDecoder(reader)
// 	//svgBytes, err := OpenFileAsBytes(filename)
// 	//if err != nil {
// 	//	return err
// 	//}
	
// 	rects := make([]xml.Rect,2)
// 	decoder.Decode(&rects)
// 	//xml.Unmarshal(svgBytes,&rects)
// 	/*{
// 		f := strconv.ParseFloat
// 		for i, item := range svgData.Elements {
// 			item = item.(svg.Rect)
// 			entities.NewSolid(f(item.Rx, 64), f(item.R, 64),
// 				f(item.Width), f(item.Height),
// 				render.NewColorBox(f(item.Width), f(item.Height),
// 					color.RGBA{100, 100, 100}),
// 				nil, i+200)
// 		}
// 	}*/
// 	fmt.Println(rects[1])
// 	return nil
// }

// ~~level data is to be stored as json, problebly compressed in the final game~~
// obsolete, do not use
func LoadJsonLevelData(filename string,offsetX,offsetY float64) {
	dlog.Info("loading json level data from", filename)
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
		rect := entities.NewSolid(rectData.X+offsetX, rectData.Y+offsetY, rectData.W, rectData.H,
			render.NewColorBox(int(rectData.W), int(rectData.H), color.RGBA{100, 100, 100, 255}),
			nil, event.CID(i+10))

		rect.Init()
		rect.UpdateLabel(rectData.Label)
		render.Draw(rect.R)
	}
}

func LoadTmx(mapPath string) error {
	/*mapReader, err := fileutil.Open(mapPath)
	if err != nil {
		return err
	}

	gameMap, err := tiled.LoadFromReader("assets",mapReader)
	if err != nil {
		return err
	}*/
	levelMap, err := tiled.LoadFromFile(mapPath)
	if err != nil {
		return err
	}
	fmt.Println(levelMap)


	fmt.Println(levelMap.TileHeight,levelMap.TileWidth)
	LoadTileLayers(*levelMap)
	LoadObjects(*levelMap)

	return nil
}

/*func LoadTmxLayer(layer *tiled.Layer,parentMap *tiled.Map) {
	for i, tile := layer.Tiles {
		} else {
			entities.NewSolid((parentMap.TileWidth*i)%parentMap.Width, y float64, w float64, h float64, r render.Renderable, tree *collision.Tree, cid event.CID)
		}
	}
}*/

func LoadTileLayers(levelMap tiled.Map) error {
	// for each _Loop, the contents of the loop are ran once for every _
	LayerLoop: for _, layer := range levelMap.Layers {
		/* RowLoop: */ for i := 0; i < levelMap.Height; i++ {
			BlockLoop: for j := 0; j < levelMap.Width; j++ {
				tileIndex := j+(i*levelMap.Width)
				if tileIndex >= len(layer.Tiles) {
					continue LayerLoop
				}
				tile := layer.Tiles[tileIndex]
				if tile.Nil == true {
					continue BlockLoop
				} else {
					sprite, err := render.LoadSprite("assets/images","wall.png")
					if err != nil {
						return err
					}
					e := entities.NewSolid(
						float64(j*levelMap.TileWidth+layer.OffsetX),
						float64(i*levelMap.TileHeight+layer.OffsetY),
						float64(levelMap.TileWidth),
						float64(levelMap.TileHeight),
						//render.NewColorBox(levelMap.TileWidth,levelMap.TileHeight
						//color.RGBA{100,100,120,255}),
						sprite,
						nil, event.CID(100+tileIndex))
					e.Init()
					e.UpdateLabel(labels.Ground)
					_, err = render.Draw(e.R,1)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

const moduleCollectable = "modClct"

func LoadObjects(levelMap tiled.Map) {
	for _, objGroup := range levelMap.ObjectGroups {
		for _, obj := range objGroup.Objects {
			if obj.Type == moduleCollectable {
				collectables.NewModuleClct(
					obj.X + float64(objGroup.OffsetX),
					obj.Y + float64(objGroup.OffsetY),
					obj.Width, obj.Height,
					render.NewColorBox(int(obj.Width), int(obj.Height), color.RGBA{255,255,0,255}),70,obj.Name)
				//o.Init()

				
			}
		}
	}
}
