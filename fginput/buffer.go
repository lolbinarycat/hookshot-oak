package fginput

// a Buffer keeps track of the last len(buf) inputs pushed to it.
type Buffer struct {
	// buf is a slice who's index cycles, if len(buf) == 6 and
	// index == 3, the order in which the elements were added
	// (with 0 being the most recent) would look somthing
	// like this:
	// [3, 2, 1, 0, 5, 4]
	buf []Input
	// index is the position of the most recent input.
	// this is used so only one input ever has to written
	// every time Push is called.
	index int
}

func NewBuffer(len int) *Buffer {
	return &Buffer{
		buf: make([]Input,len),
		index: 0,
	}
}

func (b *Buffer) Push(inp Input) {
	b.index = (b.index + 1) % len(b.buf)
	b.buf[b.index] = inp
}

func (b *Buffer) PushDir(dir Direction) {
	b.Push(DirToInput(dir))
}

func (b *Buffer) PushN(inps []Input) {
	for _, inp := range inps {
		b.Push(inp)
	}
}

func (b *Buffer) Check(seq []Input) bool {
	last := b.GetNUnique(len(seq))
	if len(last) != len(seq) {
		return false
	}

	for i := range seq {
		// we use `len(last)-1-i` to account for the
		// LIFO nature of GetNUnique
		if seq[i] != last[len(last)-1-i] { return false }
	}
	return true
}

// BUG(binarycat): if n is greater than len(b.buf), values may be copied
// Get returns the n most recenly pushed values.
func (b *Buffer) GetN(n int) []Input {
	ret := make([]Input,n)

	for i := range ret {
		ret[i] = b.Get(i)
	}
	return ret
}

// Get accounts for the cycling of b.buf and returns the
// input pushed idx pushes ago 
func (b *Buffer) Get(idx int) Input {
	// the equation used for the index gets the
	// index of the item idx spots behind b.index,
	// wrapping arround when neccecary.
	return b.buf[(b.index+len(b.buf)-idx)%len(b.buf)]
}

// GetNUnique gets the first n non-repeating inputs.
// LIFO (Last In First Out)
// Example:
//   b := NewBuffer(4)
//   b.PushN([]Input{Up,Up,Down,Down})
//   b.GetNUnique(2) // []Input{Down,Up}
func (b *Buffer) GetNUnique(n int) []Input {
	ret := make([]Input,0,n)
	var i int
	for len(ret) < n && i < len(b.buf) {
		if len(ret) == 0 || ret[len(ret)-1] != b.Get(i) {
			ret = append(ret,b.Get(i))
		}
		i++
	}
	return ret
}
