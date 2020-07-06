package collectables

import (
	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/dlog"

	"github.com/lolbinarycat/hookshot-oak/labels"
	"github.com/lolbinarycat/hookshot-oak/player"
)

// ModuleClct is a Collectable that gives you a module
type ModuleClct struct {
	React *entities.Reactive
	Name string // The name of the module the collectable represents
}


func (m ModuleClct) ClctdBy(p *player.Player) {
	p.Mods[m.Name].Obtain()
	dlog.Info("module",m.Name,"collected")
}

func NewModuleClct(x, y, w, h float64, r render.Renderable, cID event.CID, modName string) (m ModuleClct) {
	m.React = entities.NewReactive(x,y,w,h,r,nil,cID)
	m.React.Init()
	m.Name = modName
	m.React.RSpace.Add(labels.Player,func(s1,s2 *collision.Space) {
		if s2 == player.GetPlayer(0).Body.Space {
			m.ClctdBy(player.GetPlayer(0))
			m.Destroy()
		}

	})
	m.React.Bind(func(_ int, _ interface{}) int {
		m.React.RSpace.CallOnHits()
		return 0
	}, event.Enter)
	return
}

func (m *ModuleClct) Destroy() {
	m.React.Destroy()
	m = nil
}
