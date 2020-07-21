package util

import (
	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/render"
)



type Point Pos

// Pos is an interface that will allow acting on differnt kinds of positions
type Pos interface {
	X() float64
	Y() float64
	SetPos(x, y float64)
}

type Rect interface {
	Pos
	GetDims() (float64,float64)
}

// RectPos2 is a rectangle defined by a minimum and maximum point
type RectPos2 struct {
	Pos // Min point, not called min to take advantage of type embedding
	Max Pos
}

func (r RectPos2) GetDims() (w,h float64) {
	return r.Max.X()-r.X(), r.Max.Y()-r.Y()
}

// RectRenderable wraps render.Renderable so it fufils Rect
type RectRenderable struct {
	render.Renderable
}

func (r RectRenderable) GetDims() (float64,float64) {
	return i2f_2(r.Renderable.GetDims())
}

// function i2f_2 takes 2 ints and returns the equivelent 2 floats
func i2f_2 (a,b int) (c,d float64) {
	return float64(a),float64(b)
}

func CenterPointInRect(p Point,r Rect) {
	w, h := r.GetDims()
	p.SetPos(w/2+r.X(), h/2+r.Y())
}

type PosFloatgeom floatgeom.Point2

func (p PosFloatgeom) X() float64 {
	return p[0]
}

func (p PosFloatgeom) Y() float64 {
	return p[1]
}

func (p *PosFloatgeom) SetPos(x, y float64) {
	*p = PosFloatgeom{x,y}
}
