package numeric

// 範囲集計ライブラリ

// Sumi -- int配列の合計
func Sumi(a []int) int {
	s := 0
	for _, x := range a {
		s += x
	}
	return s
}

// Aggregatei -- int配列の集計
func Aggregatei(init int, a []int, fn func(sum, next int) int) int {
	s := init
	for _, x := range a {
		s = fn(s, x)
	}
	return s
}
