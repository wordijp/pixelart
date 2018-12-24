package main

import (
	//"fmt"
	"os"
	"testing"

	//"pixela_art/src/lib/color"
	"pixela_art/src/lib/image"
)

const (
	psdfile    = testdata + "/layout/calendar.psd"
	layoutfile = tmp + "/calendar_layout.dat"
)

// TestDotParseEncodeDecode -- 画像のパース、エンコード、デコードをテストする
// NOTE: datfile作成も兼ねる
func TestPsdLayoutParseEncodeDecode(t *testing.T) {
	var layout image.LayoutData
	// TEST: パーステスト
	{
		file, err := os.Open(psdfile)
		if err != nil {
			t.Errorf("error open: %s", err)
		}
		defer file.Close()

		layout, err = image.ParsePsdLayout(file)
		if err != nil {
			t.Errorf("error parse: %s", err)
		}
	}

	// TEST: エンコードテスト
	{
		file, err := os.OpenFile(layoutfile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			t.Errorf("error create: %s", err)
		}
		defer file.Close()

		err = layout.Save(file)
		if err != nil {
			t.Errorf("error save: %s", err)
		}
	}

	var layout2 image.LayoutData
	// TEST: デコードテスト
	{
		file, err := os.Open(layoutfile)
		if err != nil {
			t.Errorf("error open: %s", err)
		}
		defer file.Close()

		layout2, err = image.LoadPsdLayoutData(file)
		if err != nil {
			t.Errorf("error load: %s", err)
		}
	}

	// TEST: エンコード、デコード後のデータの内容が同じか
	{
		if len(layout.Bg.Elems) != len(layout2.Bg.Elems) {
			t.Errorf("len(Elems) wrong(%d != %d)", len(layout.Bg.Elems), len(layout2.Bg.Elems))
		}

		// BackgroundLayer
		for i, length := 0, len(layout.Bg.Elems); i < length; i++ {
			if layout.Bg.Elems[i].X != layout2.Bg.Elems[i].X {
				t.Errorf("X wrong")
			}
			if layout.Bg.Elems[i].Y != layout2.Bg.Elems[i].Y {
				t.Errorf("Y wrong")
			}
			if layout.Bg.Elems[i].Rgb.R != layout2.Bg.Elems[i].Rgb.R {
				t.Errorf("R wrong")
			}
			if layout.Bg.Elems[i].Rgb.G != layout2.Bg.Elems[i].Rgb.G {
				t.Errorf("G wrong")
			}
			if layout.Bg.Elems[i].Rgb.B != layout2.Bg.Elems[i].Rgb.B {
				t.Errorf("B wrong")
			}
		}

		// PlaceLayer
		for i, length := 0, len(layout.Place.Elems); i < length; i++ {
			if len(layout.Place.Elems[i].XY) != len(layout2.Place.Elems[i].XY) {
				t.Errorf("len(Elems) wrong(%d != %d)", len(layout.Place.Elems[i].XY), len(layout2.Place.Elems[i].XY))
			}

			for j, elemLength := 0, len(layout.Place.Elems[i].XY); j < elemLength; j++ {
				if layout.Place.Elems[i].XY[j].X != layout2.Place.Elems[i].XY[j].X {
					t.Errorf("XY[j].X wrong")
				}
				if layout.Place.Elems[i].XY[j].Y != layout2.Place.Elems[i].XY[j].Y {
					t.Errorf("XY[j].Y wrong")
				}
			}

			if layout.Place.Elems[i].Rgb.R != layout2.Place.Elems[i].Rgb.R {
				t.Errorf("R wrong")
			}
			if layout.Place.Elems[i].Rgb.G != layout2.Place.Elems[i].Rgb.G {
				t.Errorf("G wrong")
			}
			if layout.Place.Elems[i].Rgb.B != layout2.Place.Elems[i].Rgb.B {
				t.Errorf("B wrong")
			}
		}

		if layout.MinX != layout2.MinX {
			t.Errorf("MinX wrong")
		}
		if layout.MaxX != layout2.MaxX {
			t.Errorf("MaxX wrong")
		}
		if layout.MinY != layout2.MinY {
			t.Errorf("MinY wrong")
		}
		if layout.MaxY != layout2.MaxY {
			t.Errorf("MaxY wrong")
		}
	}
}

// 画像を読み込んでパースと、パース後を読み込むのとどちらが速いか

// 画像を読み込み、パースする
func BenchmarkParsePsdLayoutImage(t *testing.B) {
	file, err := os.Open(psdfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		file.Seek(0, 0)

		_, err = image.ParsePsdLayout(file)
		if err != nil {
			panic(err)
		}
	}
}

// パース済み画像情報を読み込む
func BenchmarkLoadParsedPsdLayoutImage(t *testing.B) {
	file, err := os.Open(layoutfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		file.Seek(0, 0)

		_, err = image.LoadPsdLayoutData(file)
		if err != nil {
			panic(err)
		}
	}
}
