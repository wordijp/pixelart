package image

import (
	"log"
	"math/rand"

	"pixela_art/src/lib/numeric"
	"pixela_art/src/lib/svg"
	test "pixela_art/src/lib/testify"
)

// ColorLevels -- 適用する色段階一覧
type ColorLevels struct {
	levels []int
}

// BuildColorLevels -- 各ドットのcolor levelを生成する
func BuildColorLevels(dots DotImageData, agg svg.PixelaAggregateMap) ColorLevels {
	cl := ColorLevels{}
	cl.levels = make([]int, len(dots.Elems), len(dots.Elems))

	// 色段階の配列を作成
	// PixelaのSVGは365 - 371、これをドット数にスケーリングする
	dotCount := len(dots.Elems)

	aggCounts := func() []int {
		size := agg.Size()

		a := make([]int, size, size)
		for i := 0; i < size; i++ {
			_, valref := agg.NthRef(i)
			a[i] = *valref
		}
		return a
	}()

	scaleAggCounts, _ := scaleArray(aggCounts, dotCount)

	// ドット分の色levelをセット
	// とりあえず5段階だけ(無色 + プラス4段)
	// TODO: Pixelaではマイナスあもある(9段階(マイナス4段 + 無色 + プラス4段)
	levelI := 0
LEVEL_FOR:
	for i, count := range scaleAggCounts {
		level := i
		for j := 0; j < count; j++ {
			cl.levels[levelI] = level
			levelI++
			if levelI > len(cl.levels) {
				break LEVEL_FOR
			}
		}
	}
	test.Assert(levelI == len(cl.levels))

	// levelsをシャッフル
	// TODO: 各levelが良い感じに散るように(二次元の一様分布、ドット飛びがあるので難しそう)
	r := rand.New(rand.NewSource(1234)) // NOTE: 常に同じ結果にする
	r.Shuffle(len(cl.levels), func(i, j int) {
		cl.levels[i], cl.levels[j] = cl.levels[j], cl.levels[i]
	})

	return cl
}

// 配列aの合計がtotalになるようにスケールする
// @return (スケール後の配列, スケール値)
func scaleArray(a []int, total int) (scaleA []int, scale float32) {
	sum := numeric.Sumi(a)
	if total == 0 || sum == 0 {
		return
	}

	scaleA = make([]int, len(a), len(a))

	add := 0
	for {
		scale = float32(total+add) / float32(sum)

		scaleTotal := 0
		for i, x := range a {
			scaleA[i] = int(float32(x) * scale)
			scaleTotal += scaleA[i]
		}

		if scaleTotal > total {
			// ここに来る？
			log.Printf("warn: scaleTotal(%d) > total(%d)", scaleTotal, total)
		}
		if scaleTotal >= total {
			break
		}

		add += total - scaleTotal
	}

	return
}
