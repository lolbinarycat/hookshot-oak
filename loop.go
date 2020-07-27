package main

import (
	"fmt"
	"os"

	"github.com/lolbinarycat/hookshot-oak/labels"
	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"

	"github.com/lolbinarycat/hookshot-oak/camera"
	"github.com/lolbinarycat/hookshot-oak/player"
	"github.com/lolbinarycat/hookshot-oak/replay"
	"github.com/lolbinarycat/hookshot-oak/ui"
)

var Paused = false

const PauseButton = key.Enter
const ConfirmButton = key.Z // activates the current selection
const CycleButton = key.Tab // cycles the current selection

var MainSceneLoop func() bool

//type PauseScreen struct {
//	Text *render.Text
//}

func buildMainSceneFuncs() (MainSceneStart func(string, interface{}), MainSceneLoop func() bool, MainSceneEnd func() (string, *scene.Result)) {

	// pauseMenuR is used to refernce the pause menu in order to undraw it.
	// this should be it's only use
	var pauseMenuR render.Renderable
	var plr = new(player.Player)

	pauseMenu := buildPauseScreen(map[string]ui.BtnAction{
		"resume": func() {
			Paused = false
			// pauseMenuR.Undraw()
		},
		"quit": func() {
			os.Exit(0)
		},
	})

	MainSceneStart = func(_ string, _ interface{}) {
		//*plr = new(player.Player)
		plr = loadScene()
		camera.StartCameraLoop(player.GetPlayer(0).Body)
		//pauseScreen := PauseScreen{
		//	Text: render.NewStrText("Paused", 0, 0),
		//}
		//pauseScreen.Text.Center()
		//render.Draw(pauseScreen.Text, 3)

		event.GlobalBind(func(_ int, _ interface{}) int {
			Paused = !Paused
			if Paused { // executed once each time the game is paused
				var err error
				pauseMenuR, err = render.Draw(pauseMenu.GetR(), 3, 3)
				if err != nil {
					panic(err)
				}
				fmt.Println("font:", render.DefFont())
				// runtime.Breakpoint()
			} else { // executed once each time the game is unpaused
				pauseMenuR.Undraw()
			}
			return 0
		}, key.Down+PauseButton)
		event.GlobalBind(func(_ int, _ interface{}) int {
			if Paused {
				err := pauseMenu.Do(ui.CycleSelection)
				if err != nil {
					panic(err)
				}
				(pauseMenuR).Undraw()
				pauseMenuR, err = render.Draw(pauseMenu.GetR(), 3, 3)
				if err != nil {
					panic(err)
				}
			}
			return 0
		}, key.Down+CycleButton)
		event.GlobalBind(func(_ int, _ interface{}) int {
			if Paused {
				err := pauseMenu.Do(ui.Do, ui.RunAction)
				if err != nil {
					panic(err)
				}
			}
			return 0
		}, key.Down+ConfirmButton)


		// set plr.HeldDir
		event.BindPriority(
			func(_ int, _ interface{}) int {
				if replay.Active {
					plr.HeldDir = replay.CurrentDir
				} else {
					plr.HeldDir = plr.Ctrls.GetDir()
				}
				return 0
			},
			event.BindingOption{
				Event: event.Event{Name: event.Enter, CallerID: int(plr.Body.CID)},
				Priority: 32,
			},
		)

	}

	MainSceneLoop = func() bool {
		hsOffX := float64(PlayerWidth/2 - HsWidth/2)
		hsOffY := float64(PlayerHeight/2 - HsHeight/2)
		if Paused == false {
			if oak.IsDown(key.Q) {
				if oak.IsDown(key.I) {
					fmt.Println(plr)
				}
				os.Exit(0)
			}

			if plr.Body.HitLabel(labels.Checkpoint) != nil {
				plr.RespawnPos = player.Pos{X: plr.Body.X(), Y: plr.Body.Y()}
			}
			if plr.Body.HitLabel(labels.Death) != nil {
				plr.Die()
			}

			//blocks := collision.WithLabels(labels.Block)
			for _, block := range blocks {
				block.DoCollision(block.BlockUpdater(plr))
				//	block.CID.E().(PhysObject).DoCollision(block.BlockUpdater)
			}

			plr.DoCollision(plr.DoStateLoop)

			if !plr.Hs.Active {
				plr.Hs.Body.SetPos(plr.Body.X()+hsOffX, //+player.Hs.X,
					plr.Body.Y()+hsOffY) //+player.Hs.Y)
			}

			plr.Hs.DoCollision(plr.HsUpdater)

			//player.Eyes[1].SetX(5)
		} else { // if game is paused
			// Do nothing for now, later display the pause menu
		}

		if oak.IsDown(key.S) {
			err := plr.Save("save.json")
			if err != nil {
				panic(err)
			}
		}
		if oak.IsDown(key.L) {
			err := plr.Load("save.json")
			if err != nil {
				panic(err)
			}
		}

		dlog.Verb("Input:", replay.GetInputFrom(plr))
		//if plr.Mods["quickrestart"].Active() {
		//	plr.Respawn()
		//}
		return true
	}

	MainSceneEnd = func() (string, *scene.Result) {
		return "platformer", nil
	}

	// return named return values
	return
}
