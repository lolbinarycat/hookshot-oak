package fginput

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	tbuf := NewBuffer(3)
	tbuf.Push(Up)
	tbuf.Push(Down)
	tbuf.Push(Left)
	assert.Equal(t,[]Input{Left,Down,Up},tbuf.Get(3),"bad values")
	tbuf.Push(Right)
	assert.Equal(t,[]Input{Right,Left,Down},tbuf.Get(3),"bad values after additional push")
}
