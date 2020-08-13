package modules

import (
	"time"
	"github.com/oakmound/oak/v2"
)

func isButtonPressedWithin(button string, dur time.Duration) bool {
	if k, d := oak.IsHeld(button); k && (d <= dur) {
		return true
	} else {
		return false
	}
}

func IsValidInputNum(n int) bool {
	return n >= 0 && n <= 7
}

