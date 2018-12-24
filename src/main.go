package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"

	svgo "github.com/ajstarks/svgo"

	"pixela_art/src/lib/image"
	"pixela_art/src/lib/map/slicemap"
	"pixela_art/src/lib/math"
	"pixela_art/src/lib/numeric"
	"pixela_art/src/lib/svg"
)

func main() {
	// SVGの読み込み
	file, err := os.Open("./vim-pixela.svg")
	if err != nil {
		log.Fatal(err)
		return
	}

	svgs, err := svg.ParsePixelaSvg(file)
	if err != nil {
		log.Fatal(err)
		return
	}

	var dots image.DotImageData
	// 画像の読み込み
	{
		file, err := os.Open("./image/dot/vim.png")
		if err != nil {
			log.Fatal(err)
			return
		}

		dots, err = image.ParseDotImage(file)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	// save
	{
		file, err := os.OpenFile("./_tmp/dots.dat", os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = dots.Save(file)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	// 画像の各ドットにSVGの各色段階を適用する
	// 1. SVGの色段階ごとの数を集計
	agg := aggregateSvg(svgs)
	// 2. 各色段階のカラーレベル
	cl := calcColorLevels(dots, agg)
	// 3. 各ドットの色にカラーレベルを適用
	applyDots := dotsApplyColorLevels(dots, cl)

	// 画像のドットデータを、SVGに変換する

	buf := bytes.NewBuffer(nil)
	writeSvgFromDots(applyDots, buf)

	fmt.Println(buf.String())
}

// PixelaAggregateMap -- Pixela SVG集計データ
type PixelaAggregateMap = slicemap.MapStringInt

func aggregateSvg(svgs svg.PixelaData) PixelaAggregateMap {
	sort.SliceStable(svgs.Elems, func(i, j int) bool {
		return svgs.Elems[i].Count < svgs.Elems[j].Count
	})

	// color毎のcount回数を集計する
	m := PixelaAggregateMap{}
	for _, x := range svgs.Elems {
		color := x.Rgb.ToColorCode()

		idx, ok := m.Insert(color, 1)
		if !ok {
			_, val := m.NthRef(idx)
			*val++
		}
	}

	for i, size := 0, m.Size(); i < size; i++ {
		//key, val := m.NthRef(i)
		//fmt.Printf("color(%s): count(%d)\n", key, *val)
	}

	return m
}

// ColorLevels -- 適用する色段階一覧
type ColorLevels struct {
	levels []int
}

func calcColorLevels(dots image.DotImageData, agg PixelaAggregateMap) ColorLevels {
	// 適用する色段階の配列を返す

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
	Assert(levelI == len(cl.levels))

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

// カラーレベルごとの彩度・明度割合
// level0の#eeeeee(RGB:238 238 238)がおよそ10%なので、10 - 100の等分割
var gLevelPercentage = []float64{0.1, 0.325, 0.55, 0.775, 1.0}

// dotsの各ドットにColorLevelsを適用し、返す
func dotsApplyColorLevels(dots image.DotImageData, cl ColorLevels) image.DotImageData {
	d := image.DotImageData{
		Elems: make([]image.DotImageElement, len(dots.Elems), len(dots.Elems)),
		H:     dots.H,
		W:     dots.W,
	}

	for i, length := 0, len(dots.Elems); i < length; i++ {
		// 色相固定で、彩度・明度に割合適用

		rgb := dots.Elems[i].Rgb

		// HSVで計算し
		hsv := rgb.ToHSV()
		// 彩度は0へ
		hsv.S = uint8(math.Lerpf(float64(hsv.S), 0.0, 1.0-gLevelPercentage[cl.levels[i]]))
		// 明度は100へ
		hsv.V = uint8(math.Lerpf(float64(hsv.V), 100.0, 1.0-gLevelPercentage[cl.levels[i]]))

		// RGBに戻す
		d.Elems[i].Rgb = hsv.ToRGB8()

		d.Elems[i].X = dots.Elems[i].X
		d.Elems[i].Y = dots.Elems[i].Y
	}

	return d
}

// SVG文字列を作成し、返す
func writeSvgFromDots(dots image.DotImageData, w io.Writer) {
	s := svgo.New(w)

	s.Startraw(fmt.Sprintf("viewBox=\"%d %d %d %d\"", 0, 0, dots.W*10, dots.H*10))
	for _, x := range dots.Elems {
		s.Rect(int(x.X)*10, int(x.Y)*10, 9, 9, fmt.Sprintf("fill=\"%s\"", x.Rgb.ToColorCode()))
	}
	s.End()
}
