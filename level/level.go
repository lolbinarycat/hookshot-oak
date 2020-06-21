package level

import (
	"fmt"
	"image/color"

	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/event"
)

//color constants
const (
	Ground collision.Label = iota
	NoWallJump
	Death
	Checkpoint
)

var (
	Gray color.RGBA = color.RGBA{100, 100, 100, 255}
	DullRed color.RGBA = color.RGBA{100,10,10,255}
)

func newCounter(n int) (func()(event.CID)) {
	var num = n
	return func() event.CID {
		num++
		return event.CID(num)
	}
}

func LoadDevRoom() {
	n := newCounter(100)
	fmt.Println(n())
	ground := entities.NewSolid(10, 400, 500, 20,
		render.NewColorBox(500, 20, Gray),
		nil, n())
	wall1 := entities.NewSolid(40, 200, 20, 500,
		render.NewColorBox(20, 500, Gray),
		nil, n())
	wall2 := entities.NewSolid(300, 200, 20, 500,
		render.NewColorBox(20, 500, Gray),
		nil, n())
	checkpoint := entities.NewSolid(200,350, 10,10,
		render.NewColorBox(10,10,color.RGBA{0,0,255,255}),
		nil,n())
	death := entities.NewSolid(340, 240, 50,50,
		render.NewColorBox(50,50,DullRed),
		nil, n())

	ground.UpdateLabel(Ground)
	wall1.UpdateLabel(Ground)
	wall2.UpdateLabel(Ground)
	checkpoint.UpdateLabel(Checkpoint)
	death.UpdateLabel(Death)

	render.Draw(ground.R)
	render.Draw(wall1.R, 1)
	render.Draw(wall2.R, 1)
	render.Draw(checkpoint.R, 1)
	render.Draw(death.R)
}
