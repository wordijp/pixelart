package math

// Minu8 --
func Minu8(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

// Maxu8 --
func Maxu8(a, b uint8) uint8 {
	if a < b {
		return b
	}
	return a
}

// Lerpf -- 線形補間
func Lerpf(a, b, t float64) float64 {
	//return a*(1.0-t) + b*t
	return a + (b-a)*(t)
}

// Truncf -- 絶対値で切り捨て
func Truncf(v float64) float64 {
	if v > 0.0 {
		return float64(int(v))
	}
	return -float64(int(-v))
}
