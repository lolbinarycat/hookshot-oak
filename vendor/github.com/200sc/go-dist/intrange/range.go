// Package intrange holds distributions that accept and return ints
package intrange

// Range represents the ability
// to poll a struct and return an integer,
// distributed over some range dependant
// on the implementing struct.
type Range interface {
	Poll() int
	Mult(float64) Range
	InRange(int) bool
	EnforceRange(int) int
	Percentile(float64) int
}
