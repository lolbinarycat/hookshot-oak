package ui

import (
	"fmt"

	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/render"
)

type Menu struct {
	interactables []Interactable
	drawables []Drawable
	// selIndex is the index of the selected item in interactables
	selIndex int
	x, y float64
}

type Interactable interface {
	// Do takes an int to allow multiple actions
	// Interactables with one action may ignore this
	// Do is variadic so nesting interactables can pass the remaining
	// args to sub-interactables.
	Do(...Action) error
}

type Action int

type Drawable interface {
	// GetR should update the renderable and return it
	GetR() render.Renderable
	// Pos Gets the items position
	Pos() (float64,float64)
}

const (
	Do Action = iota
	CycleSelection
	Focus
	Unfocus
	RunAction
)

func (m *Menu) GetR() render.Renderable {
	cR := render.NewCompositeR()
	for _, d := range m.drawables {
		x, y := d.Pos()
		cR.AppendOffset(d.GetR(),floatgeom.Point2{x,y})
	}
	cR.SetPos(m.x,m.y)
	return cR
}

type UnknownActionError struct {
	act Action
}

func (e UnknownActionError) Error() string {
	return fmt.Sprintf("unknow action '%d'", e.act)
}

func (m *Menu) Do(acts ...Action) error {
	switch acts[0] {
	case Do:
		m.GetActive().Do(acts[1])
	case CycleSelection:
		m.GetActive().Do(Unfocus)
		m.selIndex = (m.selIndex + 1) % len(m.interactables)
		m.GetActive().Do(Focus)
	default:
		return UnknownActionError{acts[0]}
	}
	return nil
}

func (m *Menu) GetActive() Interactable {
	return m.interactables[m.selIndex]
}

func (m *Menu) AddR(r render.Renderable) {
	m.AddD(WrapR(r))
}

func (m *Menu) AddD(d Drawable) {
	m.drawables = append(m.drawables, d)
}

func (m *Menu) AddI(i Interactable) {
	m.interactables = append(m.interactables, i)
}

func (m *Menu) AddDI(di interface{Drawable;Interactable}) {
	m.AddD(di)
	m.AddI(di)
}

func newMenu(bg render.Renderable, x, y float64) *Menu {
	m := Menu{x:x,y:y}
	if bg != nil {
		m.AddR(bg)
	}
	return &m
}

