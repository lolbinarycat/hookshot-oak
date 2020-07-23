package player

import "github.com/lolbinarycat/hookshot-oak/direction"
import "github.com/oakmound/oak/v2"
import "github.com/oakmound/oak/v2/dlog"

func (c *ControlConfig) GetDir() (dir direction.Dir) {
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
