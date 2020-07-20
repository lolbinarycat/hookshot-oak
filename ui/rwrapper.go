package ui

import "github.com/oakmound/oak/v2/render"

// RWrapper wraps around a render.Renderable so it fufils Drawable
type RWrapper struct {
	r render.Renderable
}

func (w RWrapper) GetR() render.Renderable {
	return w.r
}

func (w RWrapper) Pos() (x, y float64) {
	return w.r.X(), w.r.Y()
}

func WrapR (r render.Renderable) RWrapper {
	return RWrapper{r}
}
