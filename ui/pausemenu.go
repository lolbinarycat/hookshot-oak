package ui

import (
	"image/draw"

	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"
)

type PauseMenu struct {
	*OptionList
	Paused bool
}

func (m PauseMenu) DrawOffset(buff draw.Image, xOff, yOff float64) {
	if m.Paused {
		m.OptionList.DrawOffset(buff,xOff,yOff)
	}
}

func (m PauseMenu) Draw(buff draw.Image) {
	m.DrawOffset(buff, 0, 0)
}

func (m *PauseMenu) Pause() {
	m.Paused = true
}

func (m *PauseMenu) Unpause() {
	m.Paused = false
}

func (m *PauseMenu) TogglePause() {
m.Paused = !m.Paused
}

func NewPauseMenu(x, y float64, options []*Option,
	pause, confirm, cycle string) *PauseMenu {

	pm := new(PauseMenu)
	pm.OptionList = NewOptionList(x, y, options...)

	event.Bind(func(_ int,key interface{}) int {
		switch key.(string) {
		case pause:
			pm.TogglePause()
		case confirm:
			if pm.Paused {
				pm.ActivateSelected()
			}
		case cycle:
			if pm.Paused {
				pm.Cycle()
			}
		}
		return 0
	},key.Down,int(pm.Init()))

	return pm
}
