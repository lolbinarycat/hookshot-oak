package camera

import (
	"time"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities"
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
	for {
		<-tick.C

		//oak.SetScreen(int(player.Body.X()),int(player.Body.Y()))
		if int(playerBody.X()) < camPosX*WindowWidth {
			camPosX--
			//oak.SetScreen(WindowWidth*camPosX, 0)
		} else if int(playerBody.X()) > camPosX*WindowWidth+WindowWidth {
			camPosX++
		} else if int(playerBody.Y()) > camPosY*WindowHeight+WindowHeight {
			camPosY++
		} else if int(playerBody.Y()) < camPosY*WindowHeight {
			camPosY--
		} else {
			continue //if no camera position change occured, don't update the screen positon
		}
		oak.SetScreen(WindowWidth*camPosX, WindowHeight*camPosY)
	}
}
