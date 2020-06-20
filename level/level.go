package level

import (
	"image/color"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/render"
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

func LoadDevRoom() {
	ground := entities.NewSolid(10, 400, 500, 20,
		render.NewColorBox(500, 20, Gray),
		nil, 0)
	ground2 := entities.NewSolid(40, 200, 20, 500,
		render.NewColorBox(20, 500, Gray),
		nil, 1)
	ground3 := entities.NewSolid(300, 200, 20, 500,
		render.NewColorBox(20, 500, DullRed),
		nil, 2)
	checkpoint := entities.NewSolid(200,350, 10,10,
		render.NewColorBox(10,10,color.RGBA{0,0,255,255}),
		nil,3)

	ground.UpdateLabel(Ground)
	ground2.UpdateLabel(Ground)
	ground3.UpdateLabel(Death)
	checkpoint.UpdateLabel(Checkpoint)

	render.Draw(ground.R)
	render.Draw(ground2.R, 1)
	render.Draw(ground3.R, 1)
	render.Draw(checkpoint.R, 1)
}
