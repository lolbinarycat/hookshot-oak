package modules

import (
	"fmt"
	"time"

	"github.com/oakmound/oak/v2/dlog"
	oak "github.com/oakmound/oak/v2"
)

type CtrldMod struct {
	BasicMod
	input *Input
	inputTime time.Duration
}

func (l *List) AddCtrld(key string,il InputList,inputNum int,inpTime time.Duration)  (*List) {
	defer func() {
		r := recover()
		if r != nil {
			dlog.Error(fmt.Errorf("error %w during AddCtrld(%#v,%#v,%#v,%#v)",r,key,il,inputNum,inpTime))
		}
	} ()
	mod := CtrldMod{BasicMod:BasicMod{},
		inputTime:inpTime,
	}
	if inputNum != -1 {
		mod.input = &il[inputNum]
	}
	(*l)[key] = Module(&mod)
	return l
}

func (m CtrldMod) Active() bool {
	if m.input == nil {
		return false
	}
	if m.Equipped && !m.Obtained {
		dlog.Error("module equipped but not obtained. mod:",m)
		return false
	}
	if m.Equipped {
		if oak.IsDown(m.input.Key)  {
			return true
		}
		if oak.IsDown(m.input.Button) {
			return true
		}
	}
	return false
}

func (m CtrldMod) ActivatedWithin(dur time.Duration) bool {
	if m.input == nil {
		//dlog.Error("m.input = false",m)
		return false
	}
	if m.Active() == false {
		return false
	}
	return isButtonPressedWithin(m.input.Key,dur)
}

func (m CtrldMod) JustActivated() bool {
	return m.ActivatedWithin(m.inputTime)
}







func (m *CtrldMod) Unequip() {
	m.Equipped = false
	m.input.Bound = false
	m.input = nil
}


func (m *CtrldMod) Bind(il InputList,i int)  (ok bool)  {
	if il[i].Bound == false && i != -1 {
		m.input = &il[i]
		return true
	}
	return false
}

func (m CtrldMod) GetInputNum() int {
	panic("do not use: GetInputNum")
	// for i, inp := range m.player.Ctrls.Mod {
	// 	if m.input != nil && *m.input == inp {
	// 		return i
	// 	}
	// }
	// dlog.Info("GetInputNum failed on:",m)
	// return -1
}

func (m *CtrldMod) Obtain() {
	m.BasicMod.Obtain()
}

func (m *CtrldMod) Equip() {
	m.BasicMod.Equip()
}



func (m *CtrldMod) GetBasic() *BasicMod {
	return &m.BasicMod
}

func (m *CtrldMod) SetInput(i *Input) {
	if i.Bound == false {
		m.input = i
	}
}
