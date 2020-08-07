package main

import (
	"os"
	"github.com/oakmound/oak/v2/scene"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"

	"github.com/lolbinarycat/hookshot-oak/ui"
)

var titlescreenScene scene.Scene

// TitlescreenResult is the struct that will be passed to the next scene.
type TitlescreenResult struct {
	LoadSave bool
}

func init() {
	var optList *ui.OptionList
	var res TitlescreenResult
	var startGame bool
	titlescreenScene = scene.Scene{
		Start: func(_ string, _ interface{}) {
			optList = ui.NewOptionList(20,20,[]*ui.Option{
				{"Continue", func() {
					res.LoadSave = true
					startGame = true
				},nil},
				{"New Game", func() {
					res.LoadSave = false
					startGame = true
				},nil},
				{"Quit", func () {
					os.Exit(0)
				},nil},
			}...)
			optList.Init()
			render.Draw(optList)
			event.Bind(func(_ int, k interface{}) int {
				switch k.(string) {
				case key.Tab:
					optList.Cycle()
				case key.Enter:
					optList.ActivateSelected()
				}
				return 0
			},key.Down,int(optList.CID))
		},
		Loop: func () bool {
			return !startGame
		},
		End: func () (string,*scene.Result) {
			return mainSceneName, &scene.Result{NextSceneInput:res}
		},
	}
}
