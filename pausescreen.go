package main

import (
	"image/color"
	"os"

	"github.com/lolbinarycat/hookshot-oak/ui"
	"github.com/oakmound/oak/v2/render"
)

func buildPauseScreen(paused *bool) *ui.Menu {
	const btnW, btnH = 50, 20
	style := ui.Style{
		Button: &ui.ButtonStyle{
			DefaultBg: render.NewColorBox(btnW, btnH, color.RGBA{100,100,100,255}),
			FocusedBg: render.NewColorBox(btnW, btnH, color.RGBA{100,255,100,255}),
		},
		Menu: &ui.MenuStyle{
			Bg: render.NewColorBox(200,200, color.RGBA{0,0,255,255}),
		},
	}
	menu := style.NewMenu(20,20)
	btns := style.NewButtonsWithActions(20, 50, 50, map[string]ui.BtnAction{
		"resume": func () {
			*paused = false
		},
		"quit": func () {
			os.Exit(0)
		},
	})
	for _, b := range btns {
		menu.AddDI(b)
	}
	return menu
}
