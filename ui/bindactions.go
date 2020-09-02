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
			t.ProcessInput(acts.Btns.MenuActionBtns,key.(string))
		}
		return 0
	},key.Down,int(t.Init()))
}

// ProcessInput processes inp according to btns.
// It returns true if any action was taken, and false otherwise.
func (o *OptionList) ProcessInput(btns MenuActionBtns,inp string) bool {
	switch inp {
	case btns.Confirm:				
		o.ActivateSelected()
	case btns.Next:
		o.Cycle()
	case btns.Prev:
		o.CycleBack()
	default:
		return false
	}
	return true
}
