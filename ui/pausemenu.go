package ui

import (
	"image/draw"

	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"
)

type PauseMenu struct {
	ToggleableOptionList
}

func (m PauseMenu) DrawOffset(buff draw.Image, xOff, yOff float64) {
	if m.Active {
		m.OptionList.DrawOffset(buff,xOff,yOff)
	}
}

func (m PauseMenu) Draw(buff draw.Image) {
	m.DrawOffset(buff, 0, 0)
}

func (m *PauseMenu) Pause() {
	m.Active = true
}

func (m *PauseMenu) Unpause() {
	m.Active = false
}

func (m *PauseMenu) TogglePause() {
	m.Toggle()
}

func NewPauseMenu(x, y float64, options []*Option,
	pause, confirm, cycleFwd, cycleBack string) *PauseMenu {

	pm := new(PauseMenu)
	pm.ToggleableOptionList =
		ToggleableOptionList{OptionList:NewOptionList(x, y, options...)}

	event.Bind(func(_ int,key interface{}) int {
		switch key.(string) {
		case pause:
			pm.TogglePause()
		case confirm:
			if pm.Active {
				pm.ActivateSelected()
			}
		case cycleFwd:
			if pm.Active {
				pm.Cycle()
			}
		case cycleBack:
			if pm.Active {
				pm.Cycle()
			}
		}
		return 0
	},key.Down,int(pm.Init()))

	return pm
}
