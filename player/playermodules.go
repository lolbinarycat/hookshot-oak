package player

import (
	"time"

	"github.com/oakmound/oak/v2/key"
	//"gopkg.in/mcuadros/go-defaults"
	"github.com/lolbinarycat/hookshot-oak/player/modules"
)

const Frame = time.Second/60

//JumpInputTime describes the length of time after the jump button is pressed in which it will count as the player jumping.
//Setting this to be too high may result in multiple jumps to occur for one press of the jump button, while setting it too low may result in jumps being eaten.
const JumpInputTime time.Duration = time.Millisecond * 90

const HsInputTime time.Duration = time.Second/60 * 2

// alias for backwards compatability
type PlayerModuleList = modules.List


type BasicPlayerModule = modules.BasicMod

type CtrldPlayerModule = modules.CtrldMod


type PlayerModule = modules.Module


//whether modules should be automaticaly equipped when recived (depreciated)
var autoEquipMods bool = true

type ModInput = modules.Input


type ModInputList = modules.InputList

func (cnf *ControlConfig) DefaultMapCtrls() {
	cnf.Left  = key.LeftArrow
	cnf.Right = key.RightArrow
	cnf.Up    = key.UpArrow
	cnf.Down  = key.DownArrow
	cnf.Quit  = key.Q
	cnf.Mod = ModInputList{
		modules.NewInput(key.Z,cnf.Controller,"A"),
		modules.NewInput(key.X,cnf.Controller,"b"),
		modules.NewInput(key.LeftShift,cnf.Controller,"x"),
		modules.NewInput(key.R,cnf.Controller,"y"),
	}
}

func InitMods(p *Player) {
	p.Ctrls.DefaultMapCtrls()
	p.Mods = make(PlayerModuleList)
	p.Mods.AddBasic("walljump").
		AddBasic("blockpush").
		AddBasic("blockpull"). //still not implemented
		AddBasic("fly").
		AddBasic("groundpound").
		AddBasic("groundpoundjump").
		AddBasic("hsitemgrab").
		AddBasic("itemcarry").
		AddBasic("diaghs"). // diagonal hookshot
		AddBasic("luigi"). // "luigi mode": slippery, but faster max speed.
		AddBasic("compgm"). // complex ground movement
		AddBasic("longjump").
		AddCtrld("jump",p.Ctrls.Mod,0,JumpInputTime).
		AddCtrld("climb",p.Ctrls.Mod,2,time.Minute * 20).
		AddCtrld("hs",p.Ctrls.Mod,1,HsInputTime).
		AddCtrld("xdash",p.Ctrls.Mod,-1,HsInputTime).
		AddCtrld("quickrestart",p.Ctrls.Mod,3,Frame)
}





//func (m CtrldPlayerModule) String() string {

//}




