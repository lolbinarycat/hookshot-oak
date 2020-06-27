package camera

import (
	//"math"
	"math"
	"time"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/key"
)

const (
	WindowWidth  int = 800
	WindowHeight int = 600
)

const transitionFrameCount = 20

func StartCameraLoop(playerBody *entities.Moving) {
	camTicker := time.NewTicker(time.Millisecond * 100)
	go cameraLoop(*camTicker, playerBody)
}

func cameraLoop(tick time.Ticker, playerBody *entities.Moving) {
	camPosX := 0
	camPosY := 0
	playerWidth := playerBody.W
	playerHeight := playerBody.H

	/*var transitioning bool
	var totalTransitionDelta intgeom.Point2
	var transitionDelta intgeom.Point2*/

	for {
		<-tick.C

		camPosX = int(math.Floor(
			(playerBody.X()+playerWidth/2)/float64(WindowWidth),
		))
		camPosY = int(math.Floor(
			(playerBody.Y()+playerHeight/2)/float64(WindowHeight),
		))
		
		/*if int(playerBody.X())/WindowWidth < camPosX+1 {
		 	camPosX--
		 	//oak.SetScreen(WindowWidth*camPosX, 0)
		 } else if int(playerBody.X())/WindowWidth > camPosX+1 {
		 	camPosX++
		 } else if false && int(playerBody.Y()) > camPosY*WindowHeight+WindowHeight {
		 	camPosY++
		 } else if false && int(playerBody.Y()) < camPosY*WindowHeight {
		 	camPosY--
		 } else {
			continue //if no camera position change occured, don't update the screen positon
		}*/
		oak.SetScreen(WindowWidth*camPosX, WindowHeight*camPosY)

		//oak.SetViewportBounds(camPosX*WindowWidth,camPosX*WindowHeight,
		//	camPosX*WindowWidth+WindowHeight,camPosX*WindowHeight+WindowHeight )
		if oak.IsDown(key.RightShift) {
			oak.SetScreen(int(playerBody.X())-WindowWidth/2,0)
		}

		/*dir, ok := isOffScreen(playerBody)
		if !transitioning && ok {
			transitioning = true
			totalTransitionDelta = intgeom.Point2{oak.ScreenWidth, oak.ScreenHeight}.Mul(intgeom.Point2{dir.X(), dir.Y()})
			transitionDelta = totalTransitionDelta.DivConst(transitionFrameCount)
		}
		if transitioning {
			// disable movement
			// move camera one size towards the player
			if totalTransitionDelta.X() != 0 || totalTransitionDelta.Y() != 0 {
				oak.ShiftScreen(transitionDelta.X(), transitionDelta.Y())
				totalTransitionDelta = totalTransitionDelta.Sub(transitionDelta)
			} else {
				transitioning = false
			}

		}*/

	}
}

func isOffScreen(char *entities.Moving) (intgeom.Dir2, bool) {
	x := int(char.X())
	y := int(char.Y())
	if x > oak.ViewPos.X+oak.ScreenWidth {
		return intgeom.Right, true
	}
	if y > oak.ViewPos.Y+oak.ScreenHeight {
		return intgeom.Down, true
	}
	if x+int(char.W) < oak.ViewPos.X {
		return intgeom.Left, true
	}
	if y+int(char.H) < oak.ViewPos.Y {
		return intgeom.Up, true
	}
	return intgeom.Dir2{}, false
}
