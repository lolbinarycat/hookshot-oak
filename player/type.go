package player

import (
	"time"

	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/key"

	"github.com/lolbinarycat/hookshot-oak/direction"
	"github.com/lolbinarycat/hookshot-oak/player/condition"
	"github.com/lolbinarycat/hookshot-oak/physobj"
)

type PhysObject = physobj.PhysObject

//Player is a type representing the player
//StateInit is a variable that should be set to true when changing states
//it tells the state to initialize values like StateTimer
type Player struct {
	//Body           *entities.Moving
	//ActiColls      ActiveCollisions
	physobj.PhysObject
	State          PlayerState  `json:"-"`
	StateStartTime time.Time `json:"-"`
	Mods           PlayerModuleList
	Ctrls          ControlConfig
	RespawnPos     Pos
	Hs             Hookshot `json:"-"`
	HeldObj        *physobj.Block
	Eyes           [2]*render.Sprite `json:"-"`
	HeldDir        direction.Dir `json:"-"`
}

type Hookshot struct {
	physobj.PhysObject
	X, Y   float64
	Active bool
}

type PlayerState struct {
	Start, Loop, End PlayerStateFunc
	Map  map[condition.Condition]PlayerStateMapFunc
	MaxDuration time.Duration
	NextState *PlayerState //only used when MaxDuration is reached
}

type PlayerStateFunc func(*Player)

// if a PlayerStateMapFunc returns nil, the player's state will not change
type PlayerStateMapFunc func (p *Player) *PlayerState



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
