package main

import (
	"fmt"
	"os"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/key"
	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/oakmound/oak/v2/render"
)

var Paused = false
var PauseButtonHeld = false //stop the game from pausing/unpausing every frame
const PauseButton = key.P

var MainSceneLoop func() bool

type PauseScreen struct {
	Text *render.Text
}

func init() {

	pauseScreen := PauseScreen{
		Text:render.NewStrText("Paused",0,0),
	}
	pauseScreen.Text.Center()
	render.Draw(pauseScreen.Text,3)

	MainSceneLoop = func() bool {
		hsOffX := player.Body.W/2 - player.Hs.Body.H/2
		hsOffY := player.Body.H/2 - player.Hs.Body.H/2
	if Paused == false {
			//xdlog.SetDebugLevel(dlog.VERBOSE)
			if oak.IsDown(key.L) {
				//oak.ScreenWidth = 800
				//oak.ScreenHeight = 600
				//oak.ChangeWindow(800,600)
				//oak.MoveWindow(20, 20, 800, 600)
				//oak.SetAspectRatio(16 / 9)
			}
			if oak.IsDown(currentControls.Quit) {
				if oak.IsDown(key.I) {
					fmt.Println(player)
				}
				os.Exit(0)
			}

			if player.Body.HitLabel(labels.Checkpoint) != nil {
				player.RespawnPos = Pos{X: player.Body.X(), Y: player.Body.Y()}
			}
			if player.Body.HitLabel(labels.Death) != nil {
				player.Die()
			}

			//blocks := collision.WithLabels(labels.Block)
			for _, block := range blocks {
				block.DoCollision(block.BlockUpdater)
				//	block.CID.E().(PhysObject).DoCollision(block.BlockUpdater)
			}

			player.DoCollision(player.DoStateLoop)

			if !player.Hs.Active {
				player.Hs.Body.SetPos(player.Body.X()+hsOffX, //+player.Hs.X,
					player.Body.Y()+hsOffY) //+player.Hs.Y)
			}

			player.Hs.DoCollision(HsUpdater)

			//player.Eyes[1].SetX(5)
	} else {
		// Do nothing for now, later display the pause menu
	}

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

