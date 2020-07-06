package direction

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDir_ZeroValue(t *testing.T) {
	noDir := Dir{}

	errStr := "returned true on zero value"
	if noDir.IsLeft() {
		t.Error("IsLeft",errStr)
	}
	if noDir.IsRight() {
		t.Error("IsRight",errStr)
	}
	if noDir.IsUp() {
		t.Error("IsUp",errStr)
	}
	if noDir.IsDown() {
		t.Error("IsDown",errStr)
	}
}

func TestDir_SimpleValues(t *testing.T) {
	leftDir := Dir{-1,0}
	if leftDir.IsLeft() == false {
		t.Error("{-1,0}.IsLeft() returns false")
	}
	maxRight := MaxRight()
	if maxRight.IsRight() == false {
		t.Error("MaxRight().IsRight == false")
	}
	if MaxLeft().IsLeft() == false {
		t.Error("MaxLeft().IsLeft() == false")
	}
}

func TestDir_ComplexValues(t *testing.T) {
	testDir := Dir{-7,8}
	if testDir.IsLeft() == false{
		t.Error("testDir.IsLeft() == false")
	}
	if testDir.IsRight() {
		t.Error("testDir.IsRight() == true")
	}
	targetDir := Dir{-8,8} 
	if testDir.OrthoDiagonalize() != targetDir {
		t.Error("OrthoDiagonalize error: expected {-8,8}, got",testDir.OrthoDiagonalize())
	}
}

func TestToCoeff(t *testing.T) {
	testDir := Dir{MaxInt8,MaxInt8}
	result := ToCoeff(testDir.V)
	assert.Equal(t,float64(1),result)

	assert.Equal(t,float64(-1),ToCoeff(MaxDown().V))
}

func BenchmarkIsLeft(b *testing.B) {
	noDir := Dir{}
	for i := 0; i < b.N; i++ {
		noDir.IsLeft()
	}
}

func BenchmarkMaxRight(b *testing.B) {
	testSlice := make([]Dir,b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testSlice[i] = MaxRight()
	}
}

func BenchmarkMaxRight_SingleVar(b *testing.B) {
	testVar := Dir{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testVar = MaxRight()
	}
	if testVar != MaxRight() {
		panic("MaxRight() != MaxRight()")
	}
}

func BenchmarkMaxRight_Alternitive(b *testing.B) {
	testSlice := make([]Dir,b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testSlice[i] = Dir{H:MaxInt8}
	}
}

func BenchmarkMaxRight_Alternitive2(b *testing.B) {
	testSlice := make([]Dir,b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testSlice[i] = Dir{1,0}
	}
}

func BenchmarkSliceCreate_1(b *testing.B) {
	sliceHolderSlice := make([][]int,b.N)
	const sliceSize = 5
	for i := 0; i < b.N; i++ {
		sliceHolderSlice[i] = make([]int,sliceSize)
	}
}

func BenchmarkSliceCreate_2(b *testing.B) {
	sliceHolderSlice := make([][]int,b.N)
	const sliceCap = 5
	for i := 0; i < b.N; i++ {
		sliceHolderSlice[i] = make([]int,5,sliceCap)
	}
}

func BenchmarkMathAbsCompare(b *testing.B) {
	var val1 int8 = 5
	var val2 int8 = -5
	for i := 0; i < b.N; i++ {
		if math.Abs(float64(val1)) != math.Abs(float64(val2)) {
			panic("aaaaahhhh")
		}
	}
}

func BenchmarkNegCompare(b *testing.B) {
	var val1 int8 = 5
	var val2 int8 = -5
	for i := 0; i < b.N; i++ {
		if !((val1 == val2) || (val1 == -val2)) {
			panic("aaah")
		}
	}
}
