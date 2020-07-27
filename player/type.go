package player

import (
	"time"

	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/key"

	"github.com/lolbinarycat/hookshot-oak/direction"
)
//Player is a type representing the player
//StateInit is a variable that should be set to true when changing states
//it tells the state to initialize values like StateTimer
type Player struct {
	//Body           *entities.Moving
	//ActiColls      ActiveCollisions
	PhysObject
	State          PlayerState  `json:"-"`
	StateStartTime time.Time `json:"-"`
	Mods           PlayerModuleList
	Ctrls          ControlConfig
	RespawnPos     Pos
	Hs             Hookshot `json:"-"`
	HeldObj        *entities.Moving `json:"-"`
	Eyes           [2]*render.Sprite `json:"-"`
	HeldDir        direction.Dir `json:"-"`
}

type Hookshot struct {
	PhysObject
	X, Y   float64
	Active bool
}

type PlayerState struct {
	Start, Loop, End PlayerStateFunc
	MaxDuration time.Duration
	NextState *PlayerState //only used when MaxDuration is reached
}

type PlayerStateFunc func(*Player)

type PhysObject struct {
	Body      *entities.Moving
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

type ControlConfig struct {
	Left, Right, Up, Down, Quit string //`json:"-"`
	Mod                         ModInputList //`json:"-"`
}

var currentControls ControlConfig = ControlConfig{
	Left:  key.LeftArrow,
	Right: key.RightArrow,
	Up:    key.UpArrow,
	Down:  key.DownArrow,
	Quit:  key.Q,
}

var curCtrls = &currentControls

//type Direction uint8

// const (
//	 Left Direction = iota
//	 Right
// )

type Pos struct {
	X float64
	Y float64
}
