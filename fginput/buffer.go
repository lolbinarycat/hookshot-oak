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


func (b *Buffer) Check(seq Sequence) {
	
}

// BUG(binarycat): if n is greater than len(b.buf), values may be copied
// Get returns the n most recenly pushed values.
func (b *Buffer) Get(n int) []Input {
	ret := make([]Input,n)

	for i := range ret {
		// the equation used for the index gets the
		// index of the item i spots behind b.index,
		// wrapping arround when neccecary.
		ret[i] = b.buf[(b.index+len(b.buf)-i)%len(b.buf)]
	}
	return ret
}
