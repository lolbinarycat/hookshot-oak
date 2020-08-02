// package renderable provides a custom renderable implementation for
// the Player type.
package renderable

import (
	"image/color"
	"image/draw"
	"image"

	"github.com/oakmound/oak/v2/render"
)

const PlayerWidth = 12
const PlayerHeight = 12

const (
	EyesYStart = PlayerHeight / 6
	EyesYEnd = EyesYStart * 4
	EyeGap = 4
	Eye1PosX = (PlayerWidth / 2) - (EyeGap / 2)
	Eye2PosX = (PlayerWidth / 2) + (EyeGap / 2)
)

type PlayerR struct {
	// we embed a renderable that is used as a base
	// the position of this renderable is the position of the
	// main renderable
	*render.Sprite
	Eyes Eyes
}

type Eyes struct {
	Color color.Color
	// Gap is the distance between the eyes
	Gap int
	Size image.Point
	// Offset describes the direction the eyes are looking.
	// {0, 0} means centered
	Offset image.Point
	bounds image.Point
}

func (e Eyes) At(x, y int) color.Color {
	//return e.Color
	switch x {
	case e.Offset.X + Eye1PosX, e.Offset.X + Eye2PosX:
		if y > EyesYStart + e.Offset.Y && y < EyesYEnd + e.Offset.Y {
			return e.Color
		}
	}
	return nil
	// // innerRect is a rect that contains the eyes.
	// var innerRect image.Rectangle
	// var innerRectBounds = image.Pt(e.Size.Y,(e.Size.X*2+e.Gap))

	// innerRect.Max = innerRectBounds
	// // this massive operation centers the rectangle
	// innerRect = innerRect.Add(image.Pt(
	// 	(e.bounds.X-innerRect.Max.X)/2,
	// 	(e.bounds.Y-innerRect.Max.X)/2,
	// ))
	// innerRect = innerRect.Add(e.Offset)

	// inPt := image.Pt(x,y)
	// if inPt.In(innerRect) {
	// 	// bRect is the rectangle of space between the eyes
	// 	var bRect image.Rectangle
	// 	bRect = innerRect
	// 	bRect.Max.X -= e.Size.X
	// 	bRect.Min.X += e.Size.X
	// 	if inPt.In(bRect) {
	// 		return nil
	// 	} else {
	// 		return e.Color
	// 	}
	// } else {
	// 	return nil
	// }


}

func (p PlayerR) DrawOffset(buff draw.Image, xOff, yOff float64) {
	//p.Sprite.DrawOffset(buff, xOff, yOff)
	var baseDim image.Point
	baseDim.X, baseDim.Y = p.Sprite.GetDims()
	p.Eyes.bounds = baseDim
	for ix:=0;ix < baseDim.X;ix++ {
		for iy:=0;iy < baseDim.Y;iy++ {
			var clr = p.Eyes.At(ix,iy)
			if clr == nil {
				clr = p.Sprite.GetRGBA().At(ix,iy)
			}
			buff.Set(
				ix+int(xOff+p.Sprite.X()),
				iy+int(yOff+p.Sprite.Y()),
				clr)
		}
	}
}

func (p PlayerR) Draw(buff draw.Image) {
	p.DrawOffset(buff, 0, 0)
}


func (p *PlayerR) SetEyeOffset(x, y int) {
	p.Eyes.Offset.X = x
	p.Eyes.Offset.Y = y
}
