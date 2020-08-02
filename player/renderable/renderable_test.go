package renderable

import (
	"image"
	"image/color"
	"testing"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

func TestScene(t *testing.T) {
	if testing.Short() == false {
		oak.AddScene("t",scene.Scene{
			Start:func(prevScene string, data interface{}) {
				pr := PlayerR{
					Sprite: render.NewColorBox(
						16, 16, color.RGBA{0,255,0,255}),
					Eyes: Eyes{
						Color: color.RGBA{0,0,255,255},
						Gap: 4,
						Size: image.Pt(1,4),
						Offset: image.Pt(-3,0),
					},
				}
				pr.SetPos(40,40)
				render.Draw(pr)
			},
			Loop:func() bool {
				return true
			},
		})
		oak.Init("t")
	}
}
