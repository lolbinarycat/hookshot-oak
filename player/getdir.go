package player

import "github.com/lolbinarycat/hookshot-oak/direction"
import "github.com/oakmound/oak/v2"

func (c *ControlConfig) GetDir() (dir direction.Dir) {
	if oak.IsDown(c.Up) {
		dir = dir.Add(direction.MaxUp())
	}
	if oak.IsDown(c.Down) {
		dir = dir.Add(direction.MaxDown())
	}
	if oak.IsDown(c.Left) {
		dir = dir.Add(direction.MaxLeft())
	}
	if oak.IsDown(c.Right) {
		dir = dir.Add(direction.MaxRight())
	}
	return
}
