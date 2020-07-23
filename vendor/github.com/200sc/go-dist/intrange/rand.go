package intrange

import (
	"math/rand"
	"time"
)

func nowRand() *rand.Rand {
	// We need to make sure we don't get the same rand source for
	// two ranges
	time.Sleep(1 * time.Millisecond)
	return rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
}
