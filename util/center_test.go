package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCenter(t *testing.T) {
	testPosFlgm  := PosFloatgeom{5,7}

	asrt := assert.New(t)

	asrt.Equal(float64(5), testPosFlgm.X(), "unexpected x value")

	testRect := RectPos2{&PosFloatgeom{},&PosFloatgeom{24,18}}
	centerTestPoint := PosFloatgeom{}
	CenterPointInRect(&centerTestPoint, testRect)
	width, _ := testRect.GetDims()
	t.Log("testRect Width:",width)
	asrt.Equal(PosFloatgeom{12,9},centerTestPoint,"wrong position for centered point")
}
