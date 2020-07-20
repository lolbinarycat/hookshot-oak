package ui

import "github.com/oakmound/oak/v2/render"

type Button struct {
	R render.Renderable
	Text *render.Text
	DefaultBg, FocusedBg *render.Sprite
	//CurrentBg **render.Sprite
	Focused bool
	X, Y float64
	Action BtnAction
}

// method Do fufils Interactable
func (b *Button) Do(a Action) error {
	switch a {
	case Activate:
		b.Action()
	case Focus:
		b.Focus()
	case Unfocus:
		b.Unfocus()
	default:
		return UnknownActionError{a}
	}
	return nil
}

// method GetR is required by interface Drawable
func (b *Button) GetR() render.Renderable {
	b.Update()
	return b.R
}

func (b *Button) Pos() (float64,float64) {
	return b.X, b.Y
}

func newButton(text string, defBg, focBg *render.Sprite, x, y float64) *Button {
	btn := new(Button)
	btn.Text = render.NewStrText(text, 5, 5)
	btn.DefaultBg, btn.FocusedBg = new(render.Sprite), new(render.Sprite)
	*btn.DefaultBg, *btn.FocusedBg = *defBg, *focBg
	btn.X, btn.Y = x, y

	btn.Focused = false

	btn.R = render.NewCompositeR(btn.DefaultBg,btn.Text)
	return btn
}

func newColoredButton(w,h int,text string,defC, focC Color,x,y float64) *Button {
	DefaultBg := render.NewColorBox(w,h,defC)
	FocusedBg := render.NewColorBox(w,h,focC)

	return newButton(text,DefaultBg,FocusedBg,x,y)
}

func (b *Button) Focus() {
	b.Focused = true
}

func (b *Button) Unfocus() {
	b.Focused = false
	b.Update()
}

func (b *Button) Update() {
	if b.Focused {
		b.R = render.NewCompositeR(b.FocusedBg,b.Text)
	} else {
		b.R = render.NewCompositeR(b.DefaultBg,b.Text)
	}
}
