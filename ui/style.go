package ui

import (
	"github.com/oakmound/oak/v2/render"
)

type Style struct {
	Button *ButtonStyle
	Menu *MenuStyle
}

type ButtonStyle struct {
	// size is determined by size of DefaultBg
	DefaultBg, FocusedBg *render.Sprite
}

type MenuStyle struct {
	Bg *render.Sprite
}

func (s Style) NewButton(text string, x, y float64) *Button {
	return newButton(text,s.Button.DefaultBg,s.Button.FocusedBg,x,y)
}

func (s Style) NewButtonsWithActions(seperaton float64,x, y float64, btnMap map[string]BtnAction) []*Button {
	i := 0
	_, btnH := s.Button.DefaultBg.GetDims()
	btnHf := float64(btnH)
	btnL := make([]*Button,len(btnMap))
	for name, action := range btnMap {
		btnL[i] = newButton(name, s.Button.DefaultBg, s.Button.FocusedBg,
			x,y+((btnHf+seperaton)*float64(i)))
		btnL[i].Action = action
		i++
	}
	return btnL
}

func (s Style) NewMenu(x, y float64) *Menu {
	return newMenu(s.Menu.Bg,x,y)
}
