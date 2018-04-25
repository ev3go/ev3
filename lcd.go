// Copyright ©2016 The ev3go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ev3

import (
	"errors"
	"image"
	"image/draw"

	"github.com/ev3go/ev3dev"
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
var LCD = ev3dev.NewFrameBuffer("/dev/fb0", NewRGBAWith, LCDWidth, LCDHeight, LCDStride)

// NewRGBAWith returns a new image.RGBA image with the given bounds
// and stride, backed by the []byte, pix. If stride is zero, a working
// stride is computed. If the length of pix is less than stride*h, an
// error is returned.
func NewRGBAWith(pix []byte, r image.Rectangle, stride int) (draw.Image, error) {
	w, h := r.Dx(), r.Dy()
	if stride == 0 {
		stride = 4 * w
	}
	if len(pix) < stride*h {
		return nil, errors.New("ev3: bad pixel buffer length")
	}
	return &image.RGBA{Pix: pix, Stride: stride, Rect: r}, nil
}
