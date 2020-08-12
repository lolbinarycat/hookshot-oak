package fginput

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSeq(t *testing.T) {
	assert.Equal(t,[]Input{Down},Seq("↓"))
	assert.Equal(t,[]Input{Down,Up},Seq("↓↑"))
	assert.Equal(t,[]Input{Down|Left},Seq("↙"))
	assert.Equal(t,[]Input{Down,Down|Right,Right},Seq("↓↘→"))
}
