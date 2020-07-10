package main

import (
	"fmt"
	"os"

	"github.com/lolbinarycat/hookshot-oak/labels"
	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/render"
	"github.com/lolbinarycat/hookshot-oak/player"
)

var Paused = false
var PauseButtonHeld = false //stop the game from pausing/unpausing every frame
const PauseButton = key.P

var MainSceneLoop func() bool

type PauseScreen struct {
	Text *render.Text
}

func initMainLoop() {

	plr := player.GetPlayer(0)
	pauseScreen := PauseScreen{
		Text: render.NewStrText("Paused", 0, 0),
	}
	pauseScreen.Text.Center()
	render.Draw(pauseScreen.Text, 3)

	MainSceneLoop = func() bool {
		//defer func () {recover()}()
		plr = player.GetPlayer(0)
		hsOffX := float64(PlayerWidth/2 - HsWidth/2)
		hsOffY := float64(PlayerHeight/2 - HsHeight/2)
		if Paused == false {
			//xdlog.SetDebugLevel(dlog.VERBOSE)
			if oak.IsDown(key.L) {
				//oak.ScreenWidth = 800
				//oak.ScreenHeight = 600
				//oak.ChangeWindow(800,600)
				//oak.MoveWindow(20, 20, 800, 600)
				//oak.SetAspectRatio(16 / 9)
			}
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
		} else {
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
		//if plr.Mods["quickrestart"].Active() {
		//	plr.Respawn()
		//}

		if oak.IsDown(PauseButton) {
			if !PauseButtonHeld {
				PauseButtonHeld = true
				Paused = !Paused
			}
		} else {
			PauseButtonHeld = false
		}
		return true
	}
}
