package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/wordijp/pixelart/graph"
	"github.com/wordijp/pixelart/layout"
)

const (
	psdfile    = image + "/layout/calendar.psd"
	graphfile  = image + "/vim-pixela.svg"
	layoutfile = tmp + "/calendar_layout.dat"
)

// TestDotParseEncodeDecode -- 画像のパース、エンコード、デコードをテストする
// NOTE: datfile作成も兼ねる
func TestPsdLayoutParseEncodeDecode(t *testing.T) {
	var layouts layout.Data
	// TEST: パーステスト
	{
		file, err := os.Open(psdfile)
		if err != nil {
			t.Errorf("error open: %s", err)
		}
		defer file.Close()

		layouts, err = layout.ParseLayoutPsd(file)
		if err != nil {
			t.Errorf("error parse: %s", err)
		}
	}

	// TEST: エンコードテスト
	{
		file, err := os.Create(layoutfile)
		if err != nil {
			t.Errorf("error create: %s", err)
		}
		defer file.Close()

		err = layouts.WriteLayoutData(file)
		if err != nil {
			t.Errorf("error save: %s", err)
		}
	}

	var layout2 layout.Data
	// TEST: デコードテスト
	{
		file, err := os.Open(layoutfile)
		if err != nil {
			t.Errorf("error open: %s", err)
		}
		defer file.Close()

		layout2, err = layout.LoadLayoutData(file)
		if err != nil {
			t.Errorf("error load: %s", err)
		}
	}

	// TEST: エンコード、デコード後のデータの内容が同じか
	{
		if len(layouts.Bg.Elems) != len(layout2.Bg.Elems) {
			t.Errorf("len(Elems) wrong(%d != %d)", len(layouts.Bg.Elems), len(layout2.Bg.Elems))
		}

		// BackgroundLayer
		for i, length := 0, len(layouts.Bg.Elems); i < length; i++ {
			if layouts.Bg.Elems[i].X != layout2.Bg.Elems[i].X {
				t.Errorf("X wrong")
			}
			if layouts.Bg.Elems[i].Y != layout2.Bg.Elems[i].Y {
				t.Errorf("Y wrong")
			}
			if layouts.Bg.Elems[i].Rgb.R != layout2.Bg.Elems[i].Rgb.R {
				t.Errorf("R wrong")
			}
			if layouts.Bg.Elems[i].Rgb.G != layout2.Bg.Elems[i].Rgb.G {
				t.Errorf("G wrong")
			}
			if layouts.Bg.Elems[i].Rgb.B != layout2.Bg.Elems[i].Rgb.B {
				t.Errorf("B wrong")
			}
		}

		// PlaceLayer
		for i, length := 0, len(layouts.Place.Elems); i < length; i++ {
			if len(layouts.Place.Elems[i].XY) != len(layout2.Place.Elems[i].XY) {
				t.Errorf("len(Elems) wrong(%d != %d)", len(layouts.Place.Elems[i].XY), len(layout2.Place.Elems[i].XY))
			}

			for j, elemLength := 0, len(layouts.Place.Elems[i].XY); j < elemLength; j++ {
				if layouts.Place.Elems[i].XY[j].X != layout2.Place.Elems[i].XY[j].X {
					t.Errorf("XY[j].X wrong")
				}
				if layouts.Place.Elems[i].XY[j].Y != layout2.Place.Elems[i].XY[j].Y {
					t.Errorf("XY[j].Y wrong")
				}
			}

			if layouts.Place.Elems[i].Rgb.R != layout2.Place.Elems[i].Rgb.R {
				t.Errorf("R wrong")
			}
			if layouts.Place.Elems[i].Rgb.G != layout2.Place.Elems[i].Rgb.G {
				t.Errorf("G wrong")
			}
			if layouts.Place.Elems[i].Rgb.B != layout2.Place.Elems[i].Rgb.B {
				t.Errorf("B wrong")
			}
		}

		if layouts.MinX != layout2.MinX {
			t.Errorf("MinX wrong")
		}
		if layouts.MaxX != layout2.MaxX {
			t.Errorf("MaxX wrong")
		}
		if layouts.MinY != layout2.MinY {
			t.Errorf("MinY wrong")
		}
		if layouts.MaxY != layout2.MaxY {
			t.Errorf("MaxY wrong")
		}
	}
}

// 画像を読み込んでパースと、パース後を読み込むのとどちらが速いか

// 画像を読み込み、パースする
func BenchmarkParseLayoutPsd(t *testing.B) {
	file, err := os.Open(psdfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		file.Seek(0, 0)

		_, err = layout.ParseLayoutPsd(file)
		if err != nil {
			panic(err)
		}
	}
}

// パース済み画像情報を読み込む
func BenchmarkLoadLayoutData(t *testing.B) {
	file, err := os.Open(layoutfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		file.Seek(0, 0)

		_, err = layout.LoadLayoutData(file)
		if err != nil {
			panic(err)
		}
	}
}

// 配置図書き出し
func BenchmarkWriteLayoutData(t *testing.B) {
	var g graph.Data
	{
		file, err := os.Open(graphfile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		g, err = graph.ParseCalendarGraphSvg(file)
		if err != nil {
			panic(err)
		}
	}

	var d layout.Data
	{
		file, err := os.Open(layoutfile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		d, err = layout.LoadLayoutData(file)
		if err != nil {
			panic(err)
		}
	}

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		d.WriteSvgString(g, ioutil.Discard)
	}
}
