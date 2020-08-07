// Package layers initializes and indexes the draw stack.
// It provides both constants describing the indexes of layers,
// and pointers to the actual layers.
package layers

import "github.com/oakmound/oak/v2/render"

// Indexes for layers
const (
	BG int = iota
	FG
	UI
)

// Actual pointers to layers
var (
	BGLayer,
	FGLayer,
	UILayer *render.RenderableHeap
)

func init() {
	BGLayer = render.NewDynamicHeap()
	FGLayer = render.NewDynamicHeap()
	UILayer = render.NewStaticHeap()

	render.SetDrawStack(BGLayer,FGLayer,UILayer)
}
