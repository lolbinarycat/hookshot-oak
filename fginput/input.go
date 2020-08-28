package fginput

func (i Input) IsLeft() bool {
	return  i&Left > 0
}

func (i Input) IsRight() bool {
	return i&Right > 0
}

func (i Input) IsUp() bool {
	return i&Up > 0
}

func (i Input) IsDown() bool {
	return i&Down > 0
}

