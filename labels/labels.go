package labels

import (
	"github.com/oakmound/oak/v2/collision"
)

const (
	Ground collision.Label = iota
	NoWallJump
	Death
	Checkpoint
)
