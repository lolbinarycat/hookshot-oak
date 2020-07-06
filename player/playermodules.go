package player

import (
	"time"
	"fmt"
	

	"github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/key"
	//"gopkg.in/mcuadros/go-defaults"
)

const Frame = time.Second/60

//JumpInputTime describes the length of time after the jump button is pressed in which it will count as the player jumping.
//Setting this to be too high may result in multiple jumps to occur for one press of the jump button, while setting it too low may result in jumps being eaten.
const JumpInputTime time.Duration = time.Millisecond * 90

const HsInputTime time.Duration = time.Second/60 * 2

type PlayerModuleList map[string]PlayerModule
	/*	Jump CtrldPlayerModule //`default`
	WallJump  BasicPlayerModule
	Climb     CtrldPlayerModule
	Hookshot  CtrldPlayerModule
	BlockPush BasicPlayerModule
	BlockPull,
	Fly,
	GroundPound, // FloorDollar
	GroundPoundJump,
	HsItemGrab BasicPlayerModule
	XDash CtrldPlayerModule
}*/

type BasicPlayerModule struct {
	Equipped bool
	Obtained bool //`default:true`
}

type CtrldPlayerModule struct {
	BasicPlayerModule
	input *ModInput
	inputTime time.Duration
}

type PlayerModule interface{
	Equip()
	Unequip()
	Obtain()
	Active() bool
	JustActivated() bool
}

// ModInput refers to a keyboard key/controller button input.
// An empty string refers to an unbound input
type ModInput struct {
	key string
	button string
}

//whether modules should be automaticaly equipped when recived (depreciated)
var autoEquipMods bool = true


func NewModInput(k string,b string) ModInput {
	return ModInput{k,b}
}

type ModInputList [8]ModInput

func (cnf *ControlConfig) DefaultMapCtrls() {
	cnf.Left  = key.LeftArrow
	cnf.Right = key.RightArrow
	cnf.Up    = key.UpArrow
	cnf.Down  = key.DownArrow
	cnf.Quit  = key.Q
	cnf.Mod = ModInputList{
		NewModInput(key.Z,""),
		NewModInput(key.X,""),
		NewModInput(key.LeftShift,""),
		NewModInput(key.R,""),
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
		AddCtrld("jump",&p.Ctrls.Mod[0],JumpInputTime).
		AddCtrld("climb",nil,time.Minute * 20).
		AddCtrld("hs",&p.Ctrls.Mod[1],HsInputTime).
		AddCtrld("xdash",&p.Ctrls.Mod[2],HsInputTime).
		AddCtrld("quickrestart",&p.Ctrls.Mod[3],Frame)
}

func (l *PlayerModuleList) AddBasic(key string) *PlayerModuleList {
	(*l)[key] = PlayerModule(&BasicPlayerModule{})
	return l
}

func (l *PlayerModuleList) AddCtrld(key string,inp *ModInput,inpTime time.Duration)  (*PlayerModuleList) {
	mod := CtrldPlayerModule{BasicPlayerModule{},inp,inpTime}
	//mod.input = inp
	//mod.inputTime = inpTime
	(*l)[key] = PlayerModule(&mod)
	return l
}

func (m CtrldPlayerModule) Active() bool {
	if m.input == nil {
return false
	}
	if m.Equipped && !m.Obtained {
		dlog.Error("module equipped but not obtained. mod:",m)
		return false
	}
	if m.Equipped {
		if oak.IsDown(m.input.key)  {
			return true
		}
		if oak.IsDown(m.input.button) {
			return true
		}
	}
	return false
}

func (m CtrldPlayerModule) ActivatedWithin(dur time.Duration) bool {
	if m.input == nil {
		//dlog.Error("m.input = false",m)
		return false
	}
	if m.Active() == false {
		return false
	}
	return isButtonPressedWithin(m.input.key,dur)
}

func (m CtrldPlayerModule) JustActivated() bool {
	return m.ActivatedWithin(m.inputTime)
}

//this is here to fufill the interface
func (m BasicPlayerModule) JustActivated() bool {
	return false
}

func (l *PlayerModuleList) GiveAll(equip bool) {
	if len(*l) == 0 {
		dlog.Error("no modules to give")
	}
	for _, m := range *l {
		m.Obtain()
		if equip {
			m.Equip()
		}
	}
}

func (m *BasicPlayerModule) Obtain() {
	m.Obtained = true
}

func (m *BasicPlayerModule) Equip() {
	if m.Obtained {
		m.Equipped = true
	}
}

func (m *BasicPlayerModule) Unequip() {
	m.Equipped = false
}

func (m *CtrldPlayerModule) Unequip() {
	m.Equipped = false
	m.input = nil
}

func (m *CtrldPlayerModule) Bind(p *Player,i int) {
	m.input = &p.Ctrls.Mod[i]
}

func (m *CtrldPlayerModule) Obtain() {
	m.BasicPlayerModule.Obtain()
}

func (m *CtrldPlayerModule) Equip() {
	m.BasicPlayerModule.Equip()
}

func (m BasicPlayerModule) Active() bool {
	return m.Equipped && m.Obtained
}



func (l PlayerModuleList) ListModules() {
	for i, m := range l {
		fmt.Println(i,m)
	}
}


func isButtonPressedWithin(button string, dur time.Duration) bool {
	if k, d := oak.IsHeld(button); k && (d <= dur) {
		return true
	} else {
		return false
	}
}
