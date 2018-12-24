package color

import (
	"fmt"
	"image/color"
	"strconv"

	"math"
)

// RGB8 -- rgbの三原色
type RGB8 struct {
	R, G, B uint8
}

// RGBA -- RGB8をrgbaに変換する
func (c RGB8) RGBA() (r, g, b, a uint32) {
	panic("do not use") // Colorインターフェースのために実装は必要だが、使わない
}

// ToHSV - RGB8 to HSV
func (c RGB8) ToHSV() HSV {
	R := float64(c.R)
	G := float64(c.G)
	B := float64(c.B)

	max := math.Max(math.Max(R, G), B)
	min := math.Min(math.Min(R, G), B)
	V := max / 255.0 * 100

	var H float64
	var S float64
	if max == min {
		H = 0.0
		S = 0.0
	} else {
		if max == R {
			H = 60.0*(G-B)/(max-min) + 0
		} else if max == G {
			H = 60.0*(B-R)/(max-min) + 120.0
		} else if max == B {
			H = 60.0*(R-G)/(max-min) + 240.0
		}

		if H > 360.0 {
			H -= 360.0
		} else if H < 0.0 {
			H += 360.0
		}
		S = (max - min) / max * 100.0
	}

	return HSV{
		uint16(H),
		uint8(S),
		uint8(V),
	}
}

var (
	// RGB8Model -- RGB8変換モデル
	RGB8Model = color.ModelFunc(rgb8Model)
)

var background = RGB8{R: 255, G: 255, B: 255} // White

func rgb8Model(c color.Color) color.Color {
	if _, ok := c.(RGB8); ok {
		return c
	}

	// 他のColor実装をRGB8へ変換する
	// NOTE: RGBAは同じuint8
	c = color.RGBAModel.Convert(c)
	rgba := c.(color.RGBA)

	alpha := float64(rgba.A) / 255.0

	return RGB8{
		R: uint8(math.Min(255.0, (1.0-alpha)*float64(background.R)+alpha*float64(rgba.R))),
		G: uint8(math.Min(255.0, (1.0-alpha)*float64(background.G)+alpha*float64(rgba.G))),
		B: uint8(math.Min(255.0, (1.0-alpha)*float64(background.B)+alpha*float64(rgba.B))),
	}
}

// ToColorCode -- HTMLカラーコードの文字列に変換する
func (c *RGB8) ToColorCode() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// RGB8FromColorCode -- HTMLカラーコードをRGB8に変換する
func RGB8FromColorCode(code string) (rgb RGB8, err error) {
	// XXX: とりあえず"#rrggbb"前提で作る

	r, err := strconv.ParseUint(code[1:3], 16, 8)
	if err != nil {
		return
	}
	g, err := strconv.ParseUint(code[3:5], 16, 8)
	if err != nil {
		return
	}
	b, err := strconv.ParseUint(code[5:7], 16, 8)
	if err != nil {
		return
	}

	rgb = RGB8{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
	}

	return
}
