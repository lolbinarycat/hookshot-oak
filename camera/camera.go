package camera

import (
	"math"
	"time"

	"github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/key"
)

const (
	WindowWidth  int = 800
	WindowHeight int = 600
)


func StartCameraLoop(playerBody *entities.Moving) {
	camTicker := time.NewTicker(time.Millisecond * 100)
	go cameraLoop(*camTicker,playerBody)
}

func cameraLoop(tick time.Ticker,playerBody *entities.Moving ) {
	camPosX := 0
	camPosY := 0
	playerWidth := playerBody.W
	//playerHeight := playerBody.H
	for {
		<-tick.C

		camPosX = int(math.Floor((playerBody.X()+playerWidth/2)/float64(WindowWidth)))
		//oak.SetScreen(int(player.Body.X()),int(player.Body.Y()))
		// if int(playerBody.X())/WindowWidth < camPosX+1 {
		// 	camPosX--
		// 	//oak.SetScreen(WindowWidth*camPosX, 0)
		// } else if int(playerBody.X())/WindowWidth > camPosX+1 {
		// 	camPosX++
		// } else if false && int(playerBody.Y()) > camPosY*WindowHeight+WindowHeight {
		// 	camPosY++
		// } else if false && int(playerBody.Y()) < camPosY*WindowHeight {
		// 	camPosY--
		// } else {
		//	continue //if no camera position change occured, don't update the screen positon
		//}
		oak.SetScreen(WindowWidth*camPosX, WindowHeight*camPosY)

		oak.SetViewportBounds(camPosX*WindowWidth,camPosX*WindowHeight,
			camPosX*WindowWidth+WindowHeight,camPosX*WindowHeight+WindowHeight )
		if oak.IsDown(key.RightShift) {
			oak.SetScreen(int(playerBody.X())-WindowWidth/2,0)
		}
	}
}
