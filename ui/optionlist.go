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
		Options:      opts,
		FocusMarker:  render.NewStrText(">", -10, 0),
	}
	return &ol
}

func (ol *OptionList) Init() event.CID {
	ol.CID = event.NextID(ol)
	return ol.CID
}

type Option struct {
	Name   string
	Action func()
	Extras *OptionExtras
}

type OptionExtras struct {
	// OnCycle is activated when the cycle button is pressed,
	// and returns whether the OptionList should cycle.
	// this allows things like SubOptionLists to capture the press and
	// use it for themselves
	OnCycle func() bool
}

func (l OptionList) String() string {
	var bldr strings.Builder
	bldr.Grow(len(l.Options) * 20)
	for _, opt := range l.Options {
		bldr.WriteString(opt.Name)
		bldr.WriteRune('\n')
	}
	return bldr.String()
}

func (l OptionList) DrawOffset(buff draw.Image, xOff, yOff float64) {
	x, y := xOff+l.X(), yOff+l.Y()
	var txt = make([]*render.Text, len(l.Options))
	txt[0] = render.NewStrText(l.Options[0].Name, x, y)
	_, lineH := txt[0].GetDims()

	for i := range txt[1:] {
		txt[i+1] = render.NewStrText(l.Options[i+1].Name, x, y)
	}

	l.FocusMarker.DrawOffset(buff, x, y+float64(lineH*l.focused))
		//y+float64(lineH*(1+l.focused)+lineH/2))
	for i, ln := range txt {
		ln.DrawOffset(buff, 0, float64(lineH*i))
	}
}

func (l OptionList) Draw(buff draw.Image) {
	l.DrawOffset(buff, 0, 0)
}

func (l *OptionList) CycleN(n int) {
	if l.Foc().Extras != nil &&
		l.Foc().Extras.OnCycle != nil {
		if l.Foc().Extras.OnCycle() == false {
			// if OnCycle return false, it means to skip the actual cycleing
			return
		}
	}
	l.focused = (l.focused + n) % len(l.Options)
}

func (l *OptionList) Cycle() {
	l.CycleN(1)
}

func (l *OptionList) CycleBack() {
	l.CycleN(-1)
}

func (l OptionList) ActivateSelected() {
	l.Options[l.focused].Action()
}

func (l OptionList) Foc() *Option {
	return l.Options[l.focused]
}
