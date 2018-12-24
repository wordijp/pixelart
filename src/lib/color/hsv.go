package color

import (
	mymath "pixela_art/src/lib/math"
)

// HSV -- hsv色空間
type HSV struct {
	H    uint16
	S, V uint8
}

// ToRGB8 - HSV to RGB8
// url) https://www.programmingalgorithms.com/algorithm/hsv-to-rgb?lang=C%2B%2B
// license) https://www.programmingalgorithms.com/terms-of-use
func (c HSV) ToRGB8() RGB8 {
	H := float64(c.H)
	S := float64(c.S) / 100.0
	V := float64(c.V) / 100.0

	r := 0.0
	g := 0.0
	b := 0.0

	if S == 0 {
		r = V
		g = V
		b = V
	} else {
		var f, p, q, t float64

		if H == 360 {
			H = 0
		} else {
			H = H / 60
		}

		i := int(mymath.Truncf(H))
		f = H - float64(i)

		p = V * (1.0 - S)
		q = V * (1.0 - (S * f))
		t = V * (1.0 - (S * (1.0 - f)))

		switch i {
		case 0:
			r = V
			g = t
			b = p
			break

		case 1:
			r = q
			g = V
			b = p
			break

		case 2:
			r = p
			g = V
			b = t
			break

		case 3:
			r = p
			g = q
			b = V
			break

		case 4:
			r = t
			g = p
			b = V
			break

		default:
			r = V
			g = p
			b = q
			break
		}

	}

	return RGB8{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
	}
}
