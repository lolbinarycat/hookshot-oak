package fginput

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBuffer(tst *testing.T) {
	asrt := assert.New(tst)
	tbuf := NewBuffer(30)
	tbuf.Push(Up)
	tbuf.Push(Down)
	tbuf.Push(Left)
	asrt.Equal([]Input{Left,Down,Up},tbuf.GetN(3),"bad values")
	tbuf.Push(Right)
	asrt.Equal([]Input{Right,Left,Down},tbuf.GetN(3),"bad values after additional push")
	tbuf.Push(Right)
	asrt.Equal([]Input{Right,Left,Down},tbuf.GetNUnique(3),"bad values from GetNUnique")

	tseq := []Input{Down,Down|Right,Right}
	tbuf.PushN([]Input{Down,Down,Down,Down|Right,Down|Right,Right})
	asrt.True(tbuf.Check(tseq),"sequence should match")
}
