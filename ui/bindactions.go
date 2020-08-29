package ui

import (
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"
)

type MenuActionBtns struct {
	Confirm, Next, Prev string
}

type ToggleMenuActionBtns struct {
	MenuActionBtns
	Pause string
}

type ToggleMenuActions struct {
	Btns ToggleMenuActionBtns
	PauseIsConfirm bool
}

func (t *ToggleableOptionList) BindActions(acts ToggleMenuActions) {	
	event.Bind(func(_ int,key interface{}) int {
		if key.(string) == acts.Btns.Pause {
			if acts.PauseIsConfirm && t.Active {
				t.ActivateSelected()
				t.Active = false
			} else {
				t.Active = !t.Active
			}
		} else if t.Active {
			switch key.(string) {
			case acts.Btns.Confirm:				
				t.ActivateSelected()
			case acts.Btns.Next:
				t.Cycle()
			case acts.Btns.Prev:
				t.CycleBack()
			}
		}
		return 0
	},key.Down,int(t.Init()))
}
