package level

import (
	"bufio"
	"fmt"
	"image/color"
	"os"

	//"strconv"


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
	// load level data
	levelMap, err := tiled.LoadFromFile(mapPath)
	if err != nil {
		return err
	}

	// create tilemap entity
	tileMapEntity := NewTileMap(float64(levelMap.TileWidth),
		float64(levelMap.TileHeight))

	err = LoadTileLayers(levelMap,tileMapEntity)
	if err != nil {
		return err
	}
	LoadObjects(levelMap)
	dlog.Info("tilemap: TileMapEntity loaded:",tileMapEntity)

	return nil
}



func LoadTileLayers(levelMap *tiled.Map, tileMapEntity *TileMap) error {
	errs := make(chan error,len(levelMap.Layers))

	for _, layer := range levelMap.Layers {
		go func (lyr *tiled.Layer) {
			errs <- LoadTileLayer(levelMap,lyr,tileMapEntity)
		} (layer)
	}

	for i:=0;i<len(levelMap.Layers);i++ {
		err := <-errs
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadTileLayer(levelMap *tiled.Map, layer *tiled.Layer,tileMapEntity *TileMap) error {
	/* RowLoop: */ for i := 0; i < levelMap.Height; i++ {
		BlockLoop: for j := 0; j < levelMap.Width; j++ {
			tileIndex := j+(i*levelMap.Width)
			if tileIndex >= len(layer.Tiles) {
				return nil
			}
			tile := layer.Tiles[tileIndex]
			if tile.Nil == true {
				continue BlockLoop
			} else {
				_, _, err := LoadTile(tile,layer,levelMap,j,i,tileMapEntity)
				if err != nil {
					return err
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
				clct := collectables.NewModuleClct(
					obj.X + float64(objGroup.OffsetX),
					obj.Y + float64(objGroup.OffsetY),
					obj.Width, obj.Height,
					render.NewColorBox(int(obj.Width), int(obj.Height), color.RGBA{255,255,0,255}),70,obj.Name)
				render.Draw(clct.React.R)
				//o.Init()
			}
		}
	}
}

// LoadTile loads `tile` from `layer` of `levelMap`.
// (x, y) is the position of the tile (in tiles), reletive to the offset of `layer`
// It returns the loaded tile (or nil), whether tile.Nil == true, and any error that occurred
// levelMap is the input map of data, while tileMap is a custom entity that is
// having tiles added to it
func LoadTile(tile *tiled.LayerTile,layer *tiled.Layer,levelMap *tiled.Map,
	x,y int, tileMap *TileMap) (
	*Tile,bool,error) {
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
		t := tileMap.AddTile(float64(x*levelMap.TileWidth+layer.OffsetX),
			float64(y*levelMap.TileWidth+layer.OffsetY),sprite)
		dlog.Verb("tilemap: Tile created:",t)
		switch tilesetTile.Type {
		case "ground":
			t.UpdateLabel(labels.Ground)
		case "spikes":
			t.UpdateLabel(labels.Death)
		case "checkpoint":
			t.UpdateLabel(labels.Checkpoint)
		case "dirt":
			t.UpdateLabel(labels.NoHs)
		case "background":
			// no label
			goto Background
		default:
			return nil, false, UnknownTileTypeError{*tilesetTile}
		}
		_, err = render.Draw(t.R,3)
	ErrCheckAndReturn:
		if err != nil {
			return nil, false, err
		}
		return t, false, nil
	Background:
		_, err = render.Draw(t.R,0,-1)
		goto ErrCheckAndReturn
	}
}

type UnknownTileTypeError struct {
	tilesetTile tiled.TilesetTile // we don't use pointers to prevent nil pointer derefrencing, we wouldn't want errors in our errors, would we?
}

func (e UnknownTileTypeError) Error() string {
	return "unknown tile type: " + e.tilesetTile.Type
}
