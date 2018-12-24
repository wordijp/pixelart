package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"pixela_art/src/lib/image"
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
	agg := svg.AggregatePixelaData(svgs)
	// 2. 各色段階のカラーレベル
	cl := image.BuildColorLevels(dots, agg)
	// 3. 各ドットの色にカラーレベルを適用
	// カラーレベルごとの彩度・明度割合
	// level0の#eeeeee(RGB:238 238 238)がおよそ10%なので、10 - 100の等分割
	levelPercentage := []float64{0.1, 0.325, 0.55, 0.775, 1.0}
	applyDots := dots.ApplyColorLevels(cl, levelPercentage)

	// 画像のドットデータを、SVGに変換する

	buf := bytes.NewBuffer(nil)
	applyDots.WriteSvg(buf)

	fmt.Println(buf.String())
}
