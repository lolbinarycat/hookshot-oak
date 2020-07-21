package ui

import (
	"image/color"
	"os"

	"github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

type Color = color.Color

type BtnAction = func ()

type Titlescreen struct {
	R render.Renderable
	Buttons []*Button
	ActiveBtn int // index of the active button
	Bg, CtrlText render.Renderable
}

func newTitlescreen(btns... *Button) Titlescreen {
	sW, sH := oak.ScreenWidth,oak.ScreenHeight
	bg := render.NewColorBox(sW,sH,color.RGBA{100,100,255,255})
	ctrlText := render.NewStrText("press Z to start",
		float64(0),float64(sH))

	ttlScrn := Titlescreen{Bg:bg,CtrlText:ctrlText,Buttons: btns}
	ttlScrn.Update()

	return ttlScrn
}

// Update updates the renderable of t
func (t *Titlescreen) Update() {
	t.R = render.NewCompositeR(t.Bg,t.CtrlText)

	for _, b := range t.Buttons {
		b.Update()
		t.R.(*render.CompositeR).AppendOffset(b.R,floatgeom.Point2{b.X,b.Y})
	}
}

func (t *Titlescreen) GetActive() *Button {
	return t.Buttons[t.ActiveBtn]
}

func (t *Titlescreen) CycleFocus() {
	t.GetActive().Unfocus()
	t.ActiveBtn = (t.ActiveBtn + 1 ) % len(t.Buttons)
	t.GetActive().Focus()
}

// AddBtn adds a button below the last button in t.Buttons, with the same attributes as that button (apart from name)
func (t *Titlescreen) AddBtn(text string, action BtnAction) {
	lBtn := t.Buttons[len(t.Buttons) - 1] // last button

	_, lBtnHeight := lBtn.DefaultBg.GetDims()

	nBtn := newButton(text, lBtn.DefaultBg, lBtn.FocusedBg, lBtn.X, lBtn.Y+float64(lBtnHeight * 2))
	nBtn.Action = action
	t.Buttons = append(t.Buttons,nBtn)
}



// BuildTitlescreenScene generates the functions for the titlescreen scene.
// It's returns are set up in such a way that they can be passed directly into oak.Add,
// without the need to store them in variables first
func BuildTitlescreenScene(thisScene, nextScene string) (
	name string, strt scene.Start,lp scene.Loop,end scene.End) {

	var startGame = new(bool)
	*startGame = false

	var ttlScrnOpts TitlescreenOptions
	var ttlScrn = new(Titlescreen)
	var cycleKeyHeld bool = false

	var ttlScrnR render.Renderable // variable for undrawing

	name = thisScene

	strt = func(_ string, _ interface{}) {
		btnDefC, btnFocC := color.RGBA{100,100,100,255}, color.RGBA{140,255,140,255}
		btnH, btnW := 20, 150
		newGameBtn := newColoredButton(btnW,btnH,"new game",btnDefC,btnFocC,100,100)
		newGameBtn.Action = func () {
			ttlScrnOpts.LoadSave = false
			*startGame = true
		}
		*ttlScrn = newTitlescreen(newGameBtn)
		ttlScrn.AddBtn("quit", func() {os.Exit(0)})
		ttlScrn.Update()
		var err error
		ttlScrnR, err = render.Draw(ttlScrn.R)
		if err != nil {
			panic(err)
		}

	}
	lp = func() bool {
		if oak.IsDown(key.Tab) {
			if cycleKeyHeld == false {
				cycleKeyHeld = true
				ttlScrn.CycleFocus()
				ttlScrn.Update()
				ttlScrnR.Undraw()
				var err error
				ttlScrnR, err = render.Draw(ttlScrn.R)
				if err != nil {
					panic(err)
				}
			}
		} else {
			cycleKeyHeld = false
		}
		if oak.IsDown(key.Z) {
			ttlScrn.GetActive().Action()
		}
		return !(*startGame)
	}
	end = func() (string,*scene.Result) {
		return nextScene, &scene.Result{NextSceneInput: ttlScrnOpts}
	}
	return
}

// TitlescreenOptions is the struct that will be passed to the next scene.
type TitlescreenOptions struct {
	LoadSave bool
}
