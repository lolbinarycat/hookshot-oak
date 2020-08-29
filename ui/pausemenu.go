package ui

import (
	"image/draw"
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

func NewPauseMenu(x, y float64, options []*Option, acts MenuActions) *PauseMenu {
	pm := new(PauseMenu)
	pm.ToggleableOptionList =
		ToggleableOptionList{OptionList:NewOptionList(x, y, options...)}

	pm.BindActions(acts)
	return pm
}
