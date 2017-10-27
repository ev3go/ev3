![Gopherbrick](gopherbrick.png)
# ev3 provides EV3-specific functions for the Go ev3dev interface

[![Build Status](https://travis-ci.org/ev3go/ev3.svg?branch=master)](https://travis-ci.org/ev3go/ev3) [![Coverage Status](https://coveralls.io/repos/ev3go/ev3/badge.svg?branch=master&service=github)](https://coveralls.io/github/ev3go/ev3?branch=master) [![GoDoc](https://godoc.org/github.com/ev3go/ev3?status.svg)](https://godoc.org/github.com/ev3go/ev3)

github.com/ev3go/ev3dev depends on ev3dev stretch. For jessie support see the [ev3dev-jessie branch](https://github.com/ev3go/ev3dev/tree/ev3dev-jessie).

## Example code

```
package main

import (
	"log"
	"time"

	"github.com/ev3go/ev3"
)

func main() {
	var bright byte
	var err error
	for i := 0; i < 10; i++ {
		err = ev3.GreenLeft.SetBrightness(int(bright)).Err()
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)

		bright = ^bright

		err = ev3.GreenRight.SetBrightness(int(bright)).Err()
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}
```

LEGO® is a trademark of the LEGO Group of companies which does not sponsor, authorize or endorse this software.
