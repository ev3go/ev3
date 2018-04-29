// Copyright ©2016 The ev3go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ev3

import (
	"github.com/ev3go/ev3dev"
	"github.com/ev3go/ev3dev/fb"
)

const (
	// LCDWidth is the width of the LCD screen in pixels.
	LCDWidth = 178

	// LCDHeight is the height of the LCD screen in pixels.
	LCDHeight = 128

	// LCDStride is the width of the LCD screen memory in bytes.
	LCDStride = 712
)

// LCD is the draw image used draw directly to the ev3 LCD screen.
// Drawing operations are safe for concurrent use, but are not atomic
// beyond the pixel level. It must be initialized before use.
var LCD = ev3dev.NewFrameBuffer("/dev/fb0", fb.NewXRGBWith, LCDWidth, LCDHeight, LCDStride)
