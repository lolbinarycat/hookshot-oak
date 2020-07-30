package physobj

import (
	"github.com/lolbinarycat/hookshot-oak/labels"
	
	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"
)

type Body = *entities.Moving

type PhysObject struct {
	Body
	ActiColls ActiveCollisions `json:"-"`
	//ExtraSolids defines labels that should be solid only for this object.
	ExtraSolids []collision.Label
}

type ActiveCollisions struct {
	GroundHit          bool
	LeftWallHit        bool
	RightWallHit       bool
	CeilingHit         bool
	HLabel, VLabel     collision.Label // these define the LAST label that was hit (horizontaly and verticaly), as ints cannot be nil
	LastHitV, LastHitH event.CID       // cid of the last collision space that was hit
}

// Block represents a pushable block
type Block struct {
	PhysObject
	Held bool
	event.CID
}

func (b *Block) Init() event.CID {
	cid := event.NextID(b)
	b.CID = cid
	//b.Body.Space.CID = cid
	//b.Body.Doodad.CID = 0
	return cid
}

// this should be set in an init function somewhere
var BlockBinding func (*Block) event.Bindable

func NewBlock(w, h float64, r render.Renderable, x, y float64) *Block {
	var blk = new(Block)
	blkId := blk.Init()
	blk.Body = entities.NewMoving(x, y, w, h, r, nil, blkId, 0)
	blk.Body.UpdateLabel(labels.Block)
	event.Bind(BlockBinding(blk),event.Enter,int(blkId))
	// next a custom event for updating whether the block is held
	event.Bind(func (_ int, data interface{}) int {
		blk.Held = data.(bool)
		return 0
	},"holdUpdate",int(blkId))

	return blk
}

