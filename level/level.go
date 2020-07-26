package level

import (
	"bufio"
	"fmt"
	"image/color"
	"os"

	//"strconv"

	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/render/mod"

	//"github.com/rustyoz/svg"
	"github.com/lafriks/go-tiled"
	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/lolbinarycat/hookshot-oak/collectables"
)


//color constants
var (
	Gray    color.RGBA = color.RGBA{100, 100, 100, 255}
	DullRed color.RGBA = color.RGBA{100, 10, 10, 255}
)


func newCounter(n int) func() event.CID {
	var num = n - 1
	return func() event.CID {
		num++
		return event.CID(num)
	}
}

var getNextCID func() event.CID

func init() {
	getNextCID = newCounter(100)
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

func LoadDevRoom() error {
	n := newCounter(3)
	fmt.Println(n())

	checkpoint := entities.NewSolid(200, 350, 10, 10,
		render.NewColorBox(10, 10, color.RGBA{0, 0, 255, 255}),
		nil, n())
	checkpoint.Init()
	death := entities.NewSolid(340, 240, 50, 50,
		render.NewColorBox(50, 50, DullRed),
		nil, n())
	death.Init()


	checkpoint.UpdateLabel(labels.Checkpoint)
	death.UpdateLabel(labels.Death)


	render.Draw(checkpoint.R, 1)
	render.Draw(death.R)



	err := LoadTmx("assets/level.tmx")
	if err != nil {
		return err
	}


	return nil
}


func LoadTmx(mapPath string) error {

	levelMap, err := tiled.LoadFromFile(mapPath)
	if err != nil {
		return err
	}

	err = LoadTileLayers(levelMap)
	if err != nil {
		return err
	}
	LoadObjects(levelMap)

	return nil
}



func LoadTileLayers(levelMap *tiled.Map) error {
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
					_, _, err := LoadTile(tile,layer,levelMap,j,i)
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

func LoadObjects(levelMap *tiled.Map) {
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

// LoadTile loads `tile` from `layer` of `levelMap`.
// (x, y) is the position of the tile (in tiles), reletive to the offset of `layer`
// It returns the loaded tile (or nil), whether tile.Nil == true, and any error that occurred
func LoadTile(tile *tiled.LayerTile,layer *tiled.Layer,levelMap *tiled.Map, x,y int) (
	*entities.Solid,bool,error) {
	if tile.Nil == true {
		return nil, true, nil
	} else {
		tilesetTile := levelMap.Tilesets[0].Tiles[tile.ID]
		spritePath := tilesetTile.Image.Source
		sprite, err := render.LoadSprite("assets/",spritePath)
		if err != nil {
			return nil, false, err
		}
		if tile.DiagonalFlip {
			sprite.SetRGBA(mod.Transpose(sprite))
		}
		if tile.HorizontalFlip {
			sprite.SetRGBA(mod.FlipX(sprite))
		}
		if tile.VerticalFlip {
			sprite.SetRGBA(mod.FlipY(sprite))
		}
		if err != nil {
			return nil, false, err
		}
		e := entities.NewSolid(
			float64(x*levelMap.TileWidth+layer.OffsetX),
			float64(y*levelMap.TileHeight+layer.OffsetY),
			float64(levelMap.TileWidth),
			float64(levelMap.TileHeight),
			sprite,
			nil, event.CID(100+x+(y*levelMap.Width)))
		e.Init()
		switch tilesetTile.Type {
		case "ground":
			e.UpdateLabel(labels.Ground)
		case "spikes":
			e.UpdateLabel(labels.Death)
		case "checkpoint":
			e.UpdateLabel(labels.Checkpoint)
		case "dirt":
			e.UpdateLabel(labels.NoHs)
		case "background":
			// no label
			goto Background
		default:
			return nil, false, UnknownTileTypeError{*tilesetTile}
		}
		_, err = render.Draw(e.R,1)
	ErrCheckAndReturn:
		if err != nil {
			return nil, false, err
		}
		return e, false, nil
	Background:
		_, err = render.Draw(e.R,0,-1)
		goto ErrCheckAndReturn
	}
}

type UnknownTileTypeError struct {
	tilesetTile tiled.TilesetTile // we don't use pointers to prevent nil pointer derefrencing, we wouldn't want errors in our errors, would we?
}

func (e UnknownTileTypeError) Error() string {
	return "unknown tile type: " + e.tilesetTile.Type
}
