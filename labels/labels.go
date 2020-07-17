package labels

import (
	"github.com/oakmound/oak/v2/collision"
)

const (
	Ground collision.Label = iota
	NoWallJump
	NoHs
	Death
	Checkpoint
	Block
	Player
	Collectable
)

var Solids []collision.Label= []collision.Label{
	Ground,
	NoWallJump,
	NoHs,
	Block,
}

var GravityAffected = []collision.Label{
	Block,
}
