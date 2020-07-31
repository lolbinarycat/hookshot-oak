package level

import (
	"fmt"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/event"
	"sync"
)

type Tile struct {
	R render.Renderable
	*collision.Space
}

func (t Tile) String() string {
	return fmt.Sprintf("{Pos:%+v,R:%v}",t.Space.Location,t.R)
}

// TileMap groups all tiles under a single CID
type TileMap struct {
	Tiles []*Tile
	event.CID
	TileWidth, TileHeight float64
	sync.Mutex
	Tree *collision.Tree
}

func NewTileMap(tw,th float64) *TileMap {
	tm := new(TileMap)
	tm.Init()
	tm.Tiles = make([]*Tile,0,1000)
	tm.TileWidth = tw
	tm.TileHeight = th
	tm.Tree = collision.DefTree
	return tm
}

func (t *TileMap) Init() event.CID {
	t.CID = event.NextID(t)
	return t.CID
}

func (t *TileMap) AddTile(x, y float64, r render.Renderable) *Tile {
	t.Lock()
	tle := new(Tile)
	//rx, ry := x*t.TileWidth, y*t.TileHeight // real x, real y
	r.SetPos(x,y)
	tle.Space = collision.NewSpace(x, y, t.TileWidth, t.TileHeight, t.CID)
	tle.R = r
	t.Tiles = append(t.Tiles,tle)
	collision.Add(tle.Space)
	t.Tree.Add(tle.Space)
	t.Unlock()
	return tle
}
