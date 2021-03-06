package player

import "github.com/lolbinarycat/hookshot-oak/direction"
import "github.com/oakmound/oak/v2"
import "github.com/oakmound/oak/v2/dlog"


// GetDir does not respect the replay system, use HeldDir instead.
func (c *ControlConfig) GetDir() (dir direction.Dir) {
	if c.Controller != nil {
		state, err := c.Controller.GetState()
		if err != nil {
			goto Keyboard
		}
		defer recover()
		return direction.Dir{H:int8(state.StickLX/256),V:int8(state.StickLY/256)}
	}
Keyboard:
	if oak.IsDown(c.Up) {
		dlog.Verb("Up held")
		dir = dir.Add(direction.MaxUp())
	}
	if oak.IsDown(c.Down) {
		dlog.Verb("Down held")
		dir = dir.Add(direction.MaxDown())
	}
	if oak.IsDown(c.Left) {
		dir = dir.Add(direction.MaxLeft())
	}
	if oak.IsDown(c.Right) {
		dir = dir.Add(direction.MaxRight())
	}
	dlog.Verb("GetDir returned",dir)
	return
}
