// Copyright ©2016 The ev3go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ev3

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var genGolden = flag.Bool("gen.golden", false, "generate golden image files")

func decodeImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

var testImages = []string{
	"gopherbrick",
	"corner",
	"black",
}

func TestMonochrome(t *testing.T) {
	for _, test := range testImages {
		golden := filepath.FromSlash("testdata/" + test + "-mono.png")

		src, err := decodeImage(filepath.FromSlash("testdata/" + test + ".png"))
		if err != nil {
			t.Fatalf("failed to read src image file %v.png: %v", test, err)
		}

		got := NewMonochrome(src.Bounds(), 0)
		draw.Draw(got, got.Bounds(), src, src.Bounds().Min, draw.Src)

		if *genGolden {
			f, err := os.Create(golden)
			if err != nil {
				t.Fatalf("failed to create golden image file %v-mono.png: %v", test, err)
			}
			defer f.Close()
			err = png.Encode(f, got)
			if err != nil {
				t.Fatalf("failed to encode golden image %v-mono.png: %v", test, err)
			}
			continue
		}

		gol, err := decodeImage(golden)
		if err != nil {
			t.Fatalf("failed to read golden image file %v-mono.png: %v", test, err)
		}
		want := NewMonochrome(gol.Bounds(), 0)
		draw.Draw(want, want.Bounds(), gol, gol.Bounds().Min, draw.Src)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Monochrome from source does not match expected image for %v test", test)
		}
	}
}

var pixelTests = []struct {
	rgb  color.RGBA
	mono Pixel
}{
	{rgb: color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}, mono: Black},
	{rgb: color.RGBA{R: 0x80, G: 0x00, B: 0x00, A: 0xff}, mono: Black},
	{rgb: color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}, mono: White},
	{rgb: color.RGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xff}, mono: Black},
	{rgb: color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}, mono: Black},
	{rgb: color.RGBA{R: 0x00, G: 0x00, B: 0x80, A: 0xff}, mono: Black},
	{rgb: color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}, mono: Black},

	{rgb: color.RGBA{R: 0x05, G: 0x0a, B: 0x0b, A: 0xff}, mono: Black},
	{rgb: color.RGBA{R: 0x0e, G: 0x21, B: 0x26, A: 0xff}, mono: Black},
	{rgb: color.RGBA{R: 0x5a, G: 0xda, B: 0xff, A: 0xff}, mono: White},
}

func (p Pixel) String() string {
	if p == Black {
		return "black"
	}
	return "white"
}

func TestMonochromeModel(t *testing.T) {
	for _, test := range pixelTests {
		got := MonochromeModel.Convert(test.rgb)
		want := test.mono
		if got != want {
			t.Errorf("unexpected Monochrome value for %+v: got: %q, want: %q", test.rgb, got, want)
		}
	}
}
