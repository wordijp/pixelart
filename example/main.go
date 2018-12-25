package main

import (
	"log"
	"os"

	"github.com/wordijp/pixelart/dot"
	"github.com/wordijp/pixelart/graph"
	"github.com/wordijp/pixelart/layout"
)

func main() {
	// Graph SVG
	var g graph.Data
	{
		file, err := os.Open("./image/vim-pixela.svg")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		g, err = graph.ParseCalendarGraphSvg(file)
	}

	doDot(g, "./image/dot/vim.png", "./output/dot-vim.svg")
	doDot(g, "./image/dot/grass.png", "./output/dot-grass.svg")

	doLayout(g, "./image/layout/calendar.psd", "./output/layout-calendar.svg")
}

func doDot(g graph.Data, srcpng, dstsvg string) {
	// ドット画像
	var dots dot.Data
	{
		file, err := os.Open(srcpng)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		dots, err = dot.ParseDotPng(file)
		if err != nil {
			log.Fatal(err)
		}
	}

	// save
	{
		file, err := os.OpenFile(dstsvg, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		dots.Convert(g).WriteSvgString(file)
	}
}

func doLayout(g graph.Data, srcpsd, dstsvg string) {
	// レイアウト画像
	var layouts layout.Data
	{
		file, err := os.Open(srcpsd)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		layouts, err = layout.ParseLayoutPsd(file)
		if err != nil {
			log.Fatal(err)
		}
	}

	// save
	{
		file, err := os.OpenFile(dstsvg, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		layouts.WriteSvgString(g, file)
	}
}
