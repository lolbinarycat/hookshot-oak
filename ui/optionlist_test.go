// +build newuitest

package ui

import (
	"os"
	"testing"

	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

func TestOptionList(t *testing.T) {
	quit := Option{
		Name: "quit",
		Action: func() {
			os.Exit(0)
		},
	}
	noop := Option{
		Name: "noop",
		Action: func() {},
	}
	noop2 := noop
	noop2.Name = "noop2"
	noop3 := noop
	noop3.Name = "noop3"
	var subopt *SubOptionList
	exitSub := Option{
		Name: "back",
		Action: func() {
			subopt.Active = false
		},
	}
	subopt = NewSubOptionList("sublist",80,45,&quit,&exitSub)
	ol := NewOptionList(30,20,&quit,&noop,subopt.BuildOption(),&noop3)
	oak.Add("test",
		func (_ string, _ interface{}) {
			render.Draw(ol)
			render.Draw(subopt)
			ol.Init()
			event.Bind(func(_ int, k interface{}) int {
				switch k.(string) {
				case key.Tab:
					ol.Cycle()
				case key.Enter:
					ol.ActivateSelected()
				}
				return 0
			},key.Down,int(ol.CID))
		},
		func () bool {
			return true
		},
		func () (string, *scene.Result) {
			return "test", nil
		})
	oak.Init("test")
}
