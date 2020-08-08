// package fginput implements functions for easier detection of fighting game inputs.
package fginput

type Input uint8

const (
	None Input = 1 << iota -1
	Up
	Down
	Left
	Right

	Invalid Input = 0xFF
)

func (i Input) Canon() Input {
	switch i {
	case None, Up, Down, Left, Right, Up&Left, Up&Right, Down&Left, Down&Right:
		return i
	case Up&Down:
		return None
	case Left&Right:
		return none
	}
}
