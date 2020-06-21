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
)

/*type Rect struct {
	
	//X, 
	//Y,
	W string `xml:"width,attr"`
	H string `xml:"height,attr"`
}*/


const (
	Ground collision.Label = iota
	NoWallJump
	Death
	Checkpoint
)

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
	var num = n
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
	n := newCounter(100)
	fmt.Println(n())
	ground := entities.NewSolid(10, 400, 500, 20,
		render.NewColorBox(500, 20, Gray),
		nil, n())
	wall1 := entities.NewSolid(40, 200, 20, 500,
		render.NewColorBox(20, 500, Gray),
		nil, n())
	wall2 := entities.NewSolid(300, 200, 20, 500,
		render.NewColorBox(20, 500, Gray),
		nil, n())
	checkpoint := entities.NewSolid(200, 350, 10, 10,
		render.NewColorBox(10, 10, color.RGBA{0, 0, 255, 255}),
		nil, n())
	death := entities.NewSolid(340, 240, 50, 50,
		render.NewColorBox(50, 50, DullRed),
		nil, n())

	ground.UpdateLabel(Ground)
	wall1.UpdateLabel(Ground)
	wall2.UpdateLabel(Ground)
	checkpoint.UpdateLabel(Checkpoint)
	death.UpdateLabel(Death)

	render.Draw(ground.R)
	render.Draw(wall1.R, 1)
	render.Draw(wall2.R, 1)
	render.Draw(checkpoint.R, 1)
	render.Draw(death.R)

	LoadJsonLevelData("level.json",-800,0)

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

//level data is to be stored as json, problebly compressed in the final game
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

		rect.UpdateLabel(rectData.Label)
		render.Draw(rect.R)
	}
}

