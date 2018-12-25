package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/wordijp/pixelart/dot"
	"github.com/wordijp/pixelart/lib/color"
)

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	teardown()
	os.Exit(ret)
}

func setup() {
	fmt.Println("setup()")
}

func teardown() {
	fmt.Println("teardown()")
}

const (
	tmp   = "../_tmp"
	image = "../example/image"
)

const (
	imgfile  = image + "/dot/vim.png"
	dotsfile = tmp + "/vim_dots.dat"
)

// TestDotParseEncodeDecode -- 画像のパース、エンコード、デコードをテストする
// NOTE: dotsfile作成も兼ねる
func TestDotParseEncodeDecode(t *testing.T) {
	var dots dot.Data
	// TEST: パーステスト
	{
		file, err := os.Open(imgfile)
		if err != nil {
			t.Errorf("error open: %s", err)
		}
		defer file.Close()

		dots, err = dot.ParseDotPng(file)
		if err != nil {
			t.Errorf("error parse: %s", err)
		}
	}

	// TEST: エンコードテスト
	{
		file, err := os.OpenFile(dotsfile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			t.Errorf("error create: %s", err)
		}
		defer file.Close()

		err = dots.WriteDotData(file)
		if err != nil {
			t.Errorf("error save: %s", err)
		}
	}

	var dots2 dot.Data
	// TEST: デコードテスト
	{
		file, err := os.Open(dotsfile)
		if err != nil {
			t.Errorf("error open: %s", err)
		}
		defer file.Close()

		dots2, err = dot.LoadDotData(file)
		if err != nil {
			t.Errorf("error load: %s", err)
		}
	}

	// TEST: エンコード、デコード後のデータの内容が同じか
	{
		if len(dots.Elems) != len(dots2.Elems) {
			t.Errorf("len(Elems) wrong(%d != %d)", len(dots.Elems), len(dots2.Elems))
		}

		length := len(dots.Elems)
		for i := 0; i < length; i++ {
			if dots.Elems[i].X != dots2.Elems[i].X {
				t.Errorf("X wrong")
			}
			if dots.Elems[i].Y != dots2.Elems[i].Y {
				t.Errorf("Y wrong")
			}
			if dots.Elems[i].Rgb.R != dots2.Elems[i].Rgb.R {
				t.Errorf("R wrong")
			}
			if dots.Elems[i].Rgb.G != dots2.Elems[i].Rgb.G {
				t.Errorf("G wrong")
			}
			if dots.Elems[i].Rgb.B != dots2.Elems[i].Rgb.B {
				t.Errorf("B wrong")
			}
		}

		if dots.MinX != dots2.MinX {
			t.Errorf("MinX wrong")
		}
		if dots.MaxX != dots2.MaxX {
			t.Errorf("MaxX wrong")
		}
		if dots.MinY != dots2.MinY {
			t.Errorf("MinY wrong")
		}
		if dots.MaxY != dots2.MaxY {
			t.Errorf("MaxY wrong")
		}
	}
}

// TestRGB8ToHSV -- RGB8 to HSVテスト
func TestRGB8ToHSV(t *testing.T) {
	rgb := color.RGB8{R: 214, G: 230, B: 133}
	hsv := rgb.ToHSV()
	if hsv.H != 69 || hsv.S != 42 || hsv.V != 90 {
		t.Errorf("HSV invalid:(%d %d %d)", hsv.H, hsv.S, hsv.V)
	}
}

// TestHSVToRGB8 -- HSV to RGB8テスト
func TestHSVToRGB8(t *testing.T) {
	hsv := color.HSV{H: 117, S: 60, V: 63}
	rgb := hsv.ToRGB8()

	if rgb.R != 69 || rgb.G != 160 || rgb.B != 64 {
		t.Errorf("RGB invalid:(%d %d %d)", rgb.R, rgb.G, rgb.B)
	}
}

// 画像を読み込んでパースと、パース後を読み込むのとどちらが速いか

// 画像を読み込み、パースする
func BenchmarkParseDotPng(t *testing.B) {
	file, err := os.Open(imgfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		file.Seek(0, 0)

		_, err = dot.ParseDotPng(file)
		if err != nil {
			panic(err)
		}
	}
}

// パース済み画像情報を読み込む
func BenchmarkLoadDotData(t *testing.B) {
	file, err := os.Open(dotsfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		file.Seek(0, 0)

		_, err = dot.LoadDotData(file)
		if err != nil {
			panic(err)
		}
	}
}
