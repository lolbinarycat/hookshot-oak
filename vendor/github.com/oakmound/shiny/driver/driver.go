// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package driver provides the default driver for accessing a screen.
package driver

// TODO: figure out what to say about the responsibility for users of this
// package to check any implicit dependencies' LICENSEs. For example, the
// driver might use third party software outside of golang.org/x, like an X11
// or OpenGL library.

import (
	"github.com/oakmound/shiny/screen"
)

// Main is called by the program's main function to run the graphical
// application.
//
// It calls f on the Screen, possibly in a separate goroutine, as some OS-
// specific libraries require being on 'the main thread'. It returns when f
// returns.
func Main(f func(screen.Screen)) {
	main(f)
}

// MonitorSize reports the size in pixels of the primary monitor.
func MonitorSize() (width int, height int) {
	return monitorSize()
}