package direction

import "math"

// type Dir represents a direction.
// it is a struct made of two componets,
// the horizontal componet, and the vertical componet.
// a value of {0, 0} represents no direction.
// {-1,0} would be purely leftward
// and {0, -1} would be purely upward.
type Dir struct{
	H, V int8
}

const MaxInt8 int8 =  127

// MinInt8 is not the actual minimum, but we use -127 for symmetry 
const MinInt8 int8 = -127

func (d Dir) IsLeft() bool {
	if d.H < 0 {
		return true
	} else {
		return false
	}
}

func (d Dir) IsRight() bool {
	if d.H > 0 {
		return true
	} else {
		return false
	}
}

func (d Dir) IsUp() bool {
	if d.V < 0 {
		return true
	} else {
		return false
	}
}

func (d Dir) IsDown() bool {
	if d.V > 0 {
		return true
	} else {
		return false
	}
}

func (d Dir) IsNothing() bool {
	return d.V == 0 && d.H == 0
}

func (d Dir) IsJustRight() bool {
	return d.IsRight() && d.V == 0
}

func (d Dir) IsJustLeft() bool {
	return d.IsLeft() && d.V == 0
}

func (d Dir) IsVert() bool {
	return d.V != 0
}

func (d Dir) IsHoriz() bool {
	return d.H != 0
}
func (d Dir) IsOrtho() bool {
	if d.H == 0 || d.V == 0 {
		return true
	} else {
		return false
	}
}

func (d Dir) IsDiag() bool {
	if d.H == d.V || -d.H == d.V {
		return true
	} else {
		return false
	}
}

func (d Dir) Orthogonalize() Dir {
	if abs(d.V) > abs(d.H) {
		return Dir{V:d.V}
	} else {
		return Dir{H:d.H}
	}
}

func (d Dir) Diagonalize() Dir {
	if abs(d.V) > abs(d.H) {
		return Dir{copySign(d.H,d.V),d.V}
	} else {
		return Dir{d.H,copySign(d.V,d.H)}
	}
}

// function OrthoDiagonalize takes the input direction 'd' and outputs the nearest
// direction that is either orthogonal or diagonal.
func (d Dir) OrthoDiagonalize() Dir {
	const threshold int8 = 64
	if abs(abs(d.H) - abs(d.V)) > 64 {
		return d.Orthogonalize()
	} else {
		return d.Diagonalize()
	}
}

func (d Dir) Maximize() Dir {
	if d.H > 0 {
		d.H = MaxInt8
	} else if d.H < 0 {
		d.H = MinInt8
	}
	if d.V > 0 {
		d.V = MaxInt8
	} else if d.V < 0 {
		d.V = MinInt8
	}

	return d
}

func MaxRight() Dir {
	return Dir{H:MaxInt8}
}

func MaxLeft() Dir {
	return Dir{H:MinInt8}
}

func MaxUp() Dir {
	return Dir{V:MinInt8}
}

func MaxDown() Dir {
	return Dir{V:MaxInt8}
}

func abs(in int8) int8 {
	return int8(math.Abs(float64(in)))
}

func sign(in int8) int8 {
	if in >= 0 {
		return 1
	} else {
		return -1
	}
}

func copySign(s,in int8) int8 {
	return sign(s)*abs(in)
}

func ToCoeff(in int8) float64 {
	return float64(in/MaxInt8)
}

func (d1 Dir) Add(d2 Dir) Dir {
	return Dir{d1.H + d2.H,d1.V + d2.V}
}

func (d Dir) HCoeff() float64 {
	return ToCoeff(d.H)
}

func (d Dir) VCoeff() float64 {
	return ToCoeff(d.V)
}
