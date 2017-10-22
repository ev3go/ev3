// Copyright ©2016 The ev3go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ev3

import (
	"fmt"

	"github.com/ev3go/ev3dev"
)

// LED handles for ev3 devices.
var (
	GreenLeft  = &ev3dev.LED{Name: led{color: "green", side: "left"}}
	GreenRight = &ev3dev.LED{Name: led{color: "green", side: "right"}}
	RedLeft    = &ev3dev.LED{Name: led{color: "red", side: "left"}}
	RedRight   = &ev3dev.LED{Name: led{color: "red", side: "right"}}
)

// led is a fmt.Stringer LED name.
type led struct {
	color string
	side  string
}

func (l led) String() string {
	var id int
	switch l.side {
	case "left":
		id = 0
	case "right":
		id = 1
	default:
		panic("ev3: invalid LED side")
	}
	return fmt.Sprintf("led%d:%s:brick-status", id, l.color)
}
