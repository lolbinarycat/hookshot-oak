// the replay package provides functions to track and play back inputs
package replay

import (
	"fmt"

	"github.com/lolbinarycat/hookshot-oak/direction"
	"github.com/lolbinarycat/hookshot-oak/player"
	"github.com/oakmound/oak/v2"

	mkey "golang.org/x/mobile/event/key"
)

// Active is whether a replay is playing back currently
var Active bool

var CurrentDir direction.Dir

type FrameInput struct {
	Dir     direction.Dir
	ModBtns uint8 // bitmask
}

func (i FrameInput) String() string {
	// we want width 2 for the hex values,
	// but for some reason I have to type 3 for it to work.
	return fmt.Sprintf("%+03X %+03X %08b", i.Dir.H, i.Dir.V, i.ModBtns)
}

func GetInputFrom(p *player.Player) FrameInput {
	var inp FrameInput
	inp.Dir = p.HeldDir
	for i, mInp := range p.Ctrls.Mod {
		if oak.IsDown(mInp.Key) { // TODO: controller support
			inp.ModBtns |= 1 << i
		}
	}
	return inp
}

func SendEventsForInput(prevInp, curInp FrameInput, ctrls player.ControlConfig) {
	if prevInp.ModBtns != curInp.ModBtns {
		var m uint8 // bitmask
		for i := 0; i < 8; i++ {
			m = 1 << i
			// pressed last frame but not this frame
			if m & prevInp.ModBtns > m & curInp.ModBtns {
				sendReleaseEvent(ctrls.Mod[i].Key)
			} else if m & prevInp.ModBtns < m & curInp.ModBtns {
				sendPressEvent(ctrls.Mod[i].Key)
			}
		}
	}
	CurrentDir = curInp.Dir
}

func sendReleaseEvent(k string) {
	ev := shinyEvents[k]
	ev.Direction = mkey.DirRelease
	oak.ShinySend(ev)
}

func sendPressEvent(k string) {
	ev := shinyEvents[k]
	ev.Direction = mkey.DirPress
	oak.ShinySend(ev)
}


var shinyEvents = map[string]mkey.Event{
	"A" : {Code:mkey.CodeA},

	"Z" : {Code:mkey.CodeZ},
}
