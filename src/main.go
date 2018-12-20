package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"pixela_art/src/lib/file"
	"pixela_art/src/lib/image"
	"pixela_art/src/lib/svg"
)

func main() {
	// SVGの読み込み
	buf, err := file.ReadAll("./vim-pixela.svg")
	if err != nil {
		log.Fatal(err)
		return
	}

	data, err := svg.ParsePixelaSvg(strings.NewReader(string(buf)))
	if err != nil {
		log.Fatal(err)
		return
	}

	// TEST: 読み込めてる？
	for _, elem := range data.Elems {
		fmt.Printf("{date: %s, count: %3d, color: %s}\n",
			elem.Date.GetString(),
			elem.Count,
			elem.Rgb.ToColorCode())
	}

	// 画像の読み込み
	file, err := os.Open("./image/dot/vim.png")
	if err != nil {
		log.Fatal(err)
		return
	}

	dots, err := image.ParseDotImage(file)
	if err != nil {
		log.Fatal(err)
		return
	}

	// save
	file, err = os.OpenFile("./dots.dat", os.O_RDWR|os.O_CREATE, 0755)
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
