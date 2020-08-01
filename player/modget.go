package player

import (
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/physics"
	oak "github.com/oakmound/oak/v2"
)

func ModGetState(prevState PlayerState,modName string) PlayerState {
	var isCtrld bool
	var mod PlayerModule
	var prevSpeed physics.Vector
	return PlayerState{
		Start: func (p *Player) {
			mod = p.Mods[modName]
			p.Mods[modName].Obtain()

			_, isCtrld = p.Mods[modName].(*CtrldPlayerModule)
			prevSpeed = p.Delta
			p.Delta.SetPos(0,0)
			dlog.Info("ModGetState: started")
		},
		Loop: func (p *Player) {
			for i, inp := range p.Ctrls.Mod {
				if inp.IsDown() {
					if isCtrld {
						if inp.Bound {
							dlog.Info("ModGetState: input",i,"already bound")
							// prevent 2 modules form being bound to the same input
							continue
						}
						dlog.Info("ModGetState: module",modName,"bound to input",i)
						mod.Equip()
						mod.(interface{SetInput(int)}).SetInput(i)
						goto Resume
					} else {
						mod.Equip()
						goto Resume
					}
				}
			}
			if oak.IsDown(key.DeleteBackspace) {
				goto Resume
			}
			return
		Resume:
			p.SetState(prevState)
			return
		},
		End: func(p *Player) {
			p.Delta = prevSpeed
		},
	}.denil()
}
