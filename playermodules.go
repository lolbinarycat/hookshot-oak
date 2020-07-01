package main

import (
	"time"
	"reflect"

	"github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/key"
	//"gopkg.in/mcuadros/go-defaults"
)

type PlayerModuleList struct {
	Jump CtrldPlayerModule //`default`
	WallJump  PlayerModule
	Climb     CtrldPlayerModule
	Hookshot  CtrldPlayerModule
	BlockPush PlayerModule
	BlockPull,
	Fly,
	GroundPound, // FloorDollar
	GroundPoundJump,
	HsItemGrab PlayerModule
}

type PlayerModule struct {
	Equipped bool
	Obtained bool //`default:true`
}

type CtrldPlayerModule struct {
	PlayerModule
	input *ModInput
	inputTime time.Duration
}

// ModInput refers to a keyboard key/controller button input.
// An empty string refers to an unbound input
type ModInput struct {
	key string
	button string
}

//whether modules should be automaticaly equipped when recived
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
	}
}

func SetDefaultCtrls(p *Player) {
	p.Ctrls.DefaultMapCtrls()

	p.Mods.Jump.input = &p.Ctrls.Mod[0]
	p.Mods.Jump.inputTime = JumpInputTime
	p.Mods.Hookshot.input = &p.Ctrls.Mod[1]
	p.Mods.Hookshot.inputTime = HsInputTime
	p.Mods.Climb.input = &p.Ctrls.Mod[2]
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
		dlog.Error("m.input = false",m)
		return false
	}
	return isButtonPressedWithin(m.input.key,dur)
}

func (m CtrldPlayerModule) JustActivated() bool {
	return m.ActivatedWithin(m.inputTime)
}

func (l *PlayerModuleList) GiveAll(equip bool) {
	rl := reflect.ValueOf(l) //reflected list
	rlv := reflect.Indirect(rl) //reflected list value
	for i := 0; i < rlv.NumField(); i++ {
		rlv.Field(i).Addr().Interface().(interface{Obtain()}).Obtain()
		if equip {
			rlv.Field(i).Addr().Interface().(interface{Equip()}).Equip()
		}
	}
}

func (m *PlayerModule) Obtain() {
	m.Obtained = true
}

func (m *PlayerModule) Equip() {
	if m.Obtained {
		m.Equipped = true
	}
}

func (m *CtrldPlayerModule) Obtain() {
	m.PlayerModule.Obtain()
}

func (m *CtrldPlayerModule) Equip() {
	m.PlayerModule.Equip()
}
