// the replay package provides functions to track and play back inputs
package replay

import (
	"fmt"

	"github.com/lolbinarycat/hookshot-oak/direction"
	"github.com/lolbinarycat/hookshot-oak/player"
	"github.com/oakmound/oak/v2"
)


type FrameInput struct {
	Dir direction.Dir
	ModBtns uint8 // bitmask
}

func (i FrameInput) String() string {
	return fmt.Sprintf("%0+2X %0+2X %08b",i.Dir.H,i.Dir.V,i.ModBtns)
}

func GetInputFrom(p *player.Player) FrameInput {
	var inp FrameInput
	inp.Dir = p.Ctrls.GetDir()
	for i, mInp := range p.Ctrls.Mod {
		if oak.IsDown(mInp.Key) { // TODO: controller support
			inp.ModBtns |= 1 << i
		}
	}
	return inp
}
