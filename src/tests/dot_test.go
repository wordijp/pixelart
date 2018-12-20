package main

import (
	"fmt"
	"os"
	"testing"

	"pixela_art/src/lib/image"
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
	tmp      = "../../_tmp"
	testdata = "../../testdata"
)

const (
	imgfile = testdata + "/dot/vim.png"
	datfile = tmp + "/vim_dots.dat"
)

// TestDotParseEncodeDecode -- 画像のパース、エンコード、デコードをテストする
func TestDotParseEncodeDecode(t *testing.T) {

	var dots image.DotImageData
	// TEST: パーステスト
	{
		file, err := os.Open(imgfile)
		if err != nil {
			t.Errorf("error open: %s", err)
		}
		defer file.Close()

		dots, err = image.ParseDotImage(file)
		if err != nil {
			t.Errorf("error parse: %s", err)
		}
	}

	// TEST: エンコードテスト
	{
		file, err := os.OpenFile(datfile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			t.Errorf("error create: %s", err)
		}
		defer file.Close()

		err = dots.Save(file)
		if err != nil {
			t.Errorf("error save: %s", err)
		}
	}

	var dots2 image.DotImageData
	// TEST: デコードテスト
	{
		file, err := os.Open(datfile)
		if err != nil {
			t.Errorf("error open: %s", err)
		}

		dots2, err = image.LoadDotImageData(file)
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
	}
}

// 画像を読み込んでパースと、パース後を読み込むのとどちらが速いか

// 画像を読み込み、パースする
func BenchmarkParseDotImage(t *testing.B) {
	for i := 0; i < t.N; i++ {
		file, err := os.Open(imgfile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		//file.Seek(0, 0)
		_, err = image.ParseDotImage(file)
		if err != nil {
			panic(err)
		}

	}
}

// パース済み画像情報を読み込む
func BenchmarkLoadParsedDotImage(t *testing.B) {
	for i := 0; i < t.N; i++ {
		file, err := os.Open(datfile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = image.LoadDotImageData(file)
		if err != nil {
			panic(err)
		}
	}
}
