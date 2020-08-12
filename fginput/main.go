// package fginput implements functions for easier detection of fighting game inputs.
package fginput

type Input uint8

const (
	None Input = 0x00
	Up Input = 1 << iota
	Down
	Left
	Right

	Invalid Input = 0xFF
)

func (i Input) Canon() Input {
	switch i {
		case None, Up, Down, Left, Right, Up|Left, Up|Right, Down|Left, Down|Right:
		return i
	case Up|Down:
		return None
	case Left|Right:
		return None
	default:
		panic("fallthrough in Input.Cannon")
	}
}


type Direction interface{
	IsLeft() bool
	IsRight() bool
	IsUp() bool
	IsDown() bool
}

func DirToInput(dir Direction) (inp Input) {
	if dir.IsLeft() {
		inp = inp|Left
	}
	if dir.IsRight() {
		inp = inp|Right
	}
	if dir.IsUp() {
		inp = inp|Up
	}
	if dir.IsDown() {
		inp = inp|Down
	}
	return inp
}

