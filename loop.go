package main

import (
	"fmt"
	"os"
	//"image"

	"github.com/lolbinarycat/hookshot-oak/labels"
	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"

	"github.com/lolbinarycat/hookshot-oak/camera"
	"github.com/lolbinarycat/hookshot-oak/player"
	//prenderable "github.com/lolbinarycat/hookshot-oak/player/renderable"
	"github.com/lolbinarycat/hookshot-oak/replay"
	"github.com/lolbinarycat/hookshot-oak/ui"
	"github.com/lolbinarycat/hookshot-oak/layers"
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
	// if nextScene is set to something other than the zero value,
	// the game will transition to that scene
	var nextScene string

	var plr = new(player.Player)

	var pauseMenu *ui.PauseMenu


	MainSceneStart = func(_ string, res interface{}) {
		pauseMenu = ui.NewPauseMenu(50,50,[]*ui.Option{
			{"Resume", func() {
				pauseMenu.Unpause()
			},nil},
			{"Titlescreen", func() {
				nextScene = "titlescreen"
			},nil},
			{"Save", func() {
				err := plr.Save("save.json")
				if err != nil {
					panic(err)
				}
			},nil},
			{"Load", func() {
				err := plr.Load("save.json")
				if err != nil {
					panic(err)
				}
			},nil},
			{"Quit", func() {
				os.Exit(0)
			},nil},
		}, PauseButton, ConfirmButton, CycleButton)
		plr = loadScene()
		pauseMenu.Active = false
		{
			_, err := render.Draw(pauseMenu,layers.UI)
			if err != nil {
				panic(err)
			}
			//pauseMenu = m.(*ui.PauseMenu)
		}

		if res.(TitlescreenResult).LoadSave {
			err := plr.Load("save.json")
			if err != nil {
				panic(err)
			}
		}
		camera.StartCameraLoop(player.GetPlayer(0).Body)

		// set plr.HeldDir
		event.BindPriority(
			func(_ int, _ interface{}) int {
				plr.LastHeldDir = plr.HeldDir
				if replay.Active {
 					plr.HeldDir = replay.CurrentDir
				} else {
					plr.HeldDir = plr.Ctrls.GetDir()
				}

				plr.R.SetEyeDir(plr.HeldDir)
				plr.DirBuffer.PushDir(plr.HeldDir)
				return 0
			},
			event.BindingOption{
				Event: event.Event{Name: event.Enter, CallerID: int(plr.Body.CID)},
				Priority: 32,
			},
		)

	}
	// 0xc0001b1720
	MainSceneLoop = func() bool {
		hsOffX := float64(PlayerWidth/2 - HsWidth/2)
		hsOffY := float64(PlayerHeight/2 - HsHeight/2)
		if pauseMenu.Active == false {
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


			plr.DoCollision(plr.DoStateLoop)

			if !plr.Hs.Active {
				plr.Hs.Body.SetPos(plr.Body.X()+hsOffX, //+player.Hs.X,
					plr.Body.Y()+hsOffY) //+player.Hs.Y)
			}

			plr.Hs.DoCollision(plr.HsUpdater)

			if (plr.ActiColls.CeilingHit && plr.ActiColls.GroundHit) ||
				(plr.ActiColls.LeftWallHit && plr.ActiColls.RightWallHit) {

				plr.Die()
			}

			//player.Eyes[1].SetX(5)
		} else { // if game is paused
			// Do nothing for now, later display the pause menu
		}


		dlog.Verb("Input:", replay.GetInputFrom(plr))
		//if plr.Mods["quickrestart"].Active() {
		//	plr.Respawn()
		//}
		return nextScene == ""
	}

	MainSceneEnd = func() (string, *scene.Result) {
		return nextScene, nil
	}

	// return named return values
	return
}
