package ui

import (
	"image/draw"
	"strings"

	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"
)

type OptionList struct {
	render.LayeredPoint
	Options []*Option
	focused int
	event.CID
	FocusMarker render.Renderable
}

func NewOptionList(x, y float64, opts ...*Option) *OptionList {
	ol := OptionList{
		LayeredPoint: render.NewLayeredPoint(x, y, 0),
		Options: opts,
		FocusMarker: render.NewStrText(">", 14, 1),
	}
	return &ol
}

func (ol *OptionList) Init() event.CID {
	ol.CID = event.NextID(ol)
	return ol.CID
}

type Option struct {
	Name string
	Action func()
}

func (l OptionList) String() string {
	var bldr strings.Builder
	bldr.Grow(len(l.Options)*20)
	for _, opt := range l.Options {
		bldr.WriteString(opt.Name)
		bldr.WriteRune('\n')
	}
	return bldr.String()
}

func (l OptionList) DrawOffset(buff draw.Image, xOff, yOff float64) {
	x, y := xOff+l.X(), yOff+l.Y()
	var txt = make([]*render.Text,len(l.Options))
	txt[0] = render.NewStrText(l.Options[0].Name, x, y)
	_ , lineH := txt[0].GetDims()

	for i := range txt[1:] {
		// through some witchcraft, y is incremented.
		txt[i+1] = render.NewStrText(l.Options[i+1].Name, x, y)
	}

	//txtW, txtH := txt[l.focused].GetDims()
	l.FocusMarker.DrawOffset(buff,x,
		y+float64(lineH*(1+l.focused)+lineH/2))
	for i, ln := range txt {
		//if ln != nil {
			ln.DrawOffset(buff, x, y+float64(lineH*i))
		//}
	}
}

func (l OptionList) Draw(buff draw.Image) {
	l.DrawOffset(buff,0,0)
}

func (l *OptionList) Cycle() {
	l.focused = (l.focused + 1) % len(l.Options)
}

func (l OptionList) ActivateSelected() {
	l.Options[l.focused].Action()
}
