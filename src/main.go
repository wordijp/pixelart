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
	defer file.Close()

	svgs, err := svg.ParsePixelaSvg(file)
	if err != nil {
		log.Fatal(err)
		return
	}

	//doDot(svgs)
	doLayout(svgs)
}

// ドットアート処理
func doDot(svgs svg.PixelaData) {
	var dots image.DotImageData
	// 画像の読み込み
	{
		file, err := os.Open("./image/dot/vim.png")
		if err != nil {
			log.Fatal(err)
			return
		}
		defer file.Close()

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
		defer file.Close()

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

// 配置アート処理
func doLayout(svgs svg.PixelaData) {
	// TODO: ...

	var layouts image.LayoutData
	// PSD画像の読み込み
	{
		file, err := os.Open("./image/layout/calendar.psd")
		if err != nil {
			log.Fatal(err)
			return
		}
		defer file.Close()

		layouts, err = image.ParsePsdLayout(file)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	//{
	//    file, err := os.Open("./_tmp/calendar_layout.dat")
	//    if err != nil {
	//        log.Fatal(err)
	//        return
	//    }
	//    defer file.Close()

	//    layouts, err = image.LoadPsdLayoutData(file)
	//}

	// save
	{
		file, err := os.OpenFile("./_tmp/layouts.dat", os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer file.Close()

		err = layouts.Save(file)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	// PSDの配置データを、SVGに変換する

	buf := bytes.NewBuffer(nil)
	layouts.WriteSvg(svgs, buf)

	fmt.Println(buf.String())
}
