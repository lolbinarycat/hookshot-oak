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
	ol := NewOptionList(30,20,&quit,&noop)
	oak.Add("test",
		func (_ string, _ interface{}) {
			render.Draw(ol)
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
