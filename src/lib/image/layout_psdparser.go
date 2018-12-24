package image

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"log"
	"math"
	"strings"
	"time"

	svgo "github.com/ajstarks/svgo"
	"github.com/oov/psd"
	"github.com/vmihailenco/msgpack"

	color "pixela_art/src/lib/color"
	"pixela_art/src/lib/svg"
	test "pixela_art/src/lib/testify"
)

// LayoutData -- 配置図用画像をパースする
type LayoutData struct {
	Bg         LayoutBackgroundLayer
	Place      LayoutPlaceLayer
	MinX, MaxX int16
	MinY, MaxY int16
}

// LayoutBackgroundLayer -- 背景用レイヤーのパースデータ
type LayoutBackgroundLayer struct {
	Elems []LayoutBackgroundElement
}

// LayoutBackgroundElement -- 背景用レイヤーのドット情報
type LayoutBackgroundElement struct {
	X, Y int16
	Rgb  color.RGB8
}

// LayoutPlaceLayer -- 配置情報用レイヤーのパースデータ
type LayoutPlaceLayer struct {
	Elems []LayoutPlaceElement
}

// LayoutPlaceElement -- 配置情報用レイヤーの配置情報
type LayoutPlaceElement struct {
	XY  []point
	Rgb color.RGB8
}

type point struct {
	X, Y int16
}

// ParsePsdLayout -- PSD画像をパースする
func ParsePsdLayout(r io.Reader) (data LayoutData, err error) {
	img, _, err := psd.Decode(r, &psd.DecodeOptions{SkipMergedImage: true})
	if err != nil {
		return
	}

	minx := math.MaxInt32
	miny := math.MaxInt32
	maxx := 0
	maxy := 0
	for _, layer := range img.Layer {
		b := layer.Rect
		if minx > b.Min.X {
			minx = b.Min.X
		}
		if miny > b.Min.Y {
			miny = b.Min.Y
		}
		if maxx < b.Max.X {
			maxx = b.Max.X
		}
		if maxy < b.Max.Y {
			maxy = b.Max.Y
		}
	}

	H := maxy - miny
	W := maxx - minx

	memo := make([]bool, H*W, H*W)
	// パース開始
	for _, layer := range img.Layer {
		if strings.Index(layer.Name, "background") >= 0 {
			bg, e := parseBackgroundLayer(layer)
			err = e
			if err != nil {
				return
			}
			data.Bg = bg
		} else if strings.Index(layer.Name, "place") >= 0 {
			// memo clear
			for i, length := 0, len(memo); i < length; i++ {
				memo[i] = false
			}

			elems, e := parsePlaceLayer(&data, layer, &memo, H, W)
			err = e
			if err != nil {
				return
			}

			data.Place.Elems = append(data.Place.Elems, elems...)
		}
	}

	data.MinX = int16(minx)
	data.MaxX = int16(maxx)
	data.MinY = int16(miny)
	data.MaxY = int16(maxy)

	return
}
func parseBackgroundLayer(bg psd.Layer) (data LayoutBackgroundLayer, err error) {
	// NOTE: 多段Layerは無視

	// 色情報をそのまま使う
	img := bg.Picker
	b := img.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			c := img.At(x, y)
			_, _, _, a := c.RGBA()
			// 透明は弾く
			if a > 0 {
				data.Elems = append(data.Elems, LayoutBackgroundElement{
					X:   int16(x),
					Y:   int16(y),
					Rgb: color.RGB8Model.Convert(c).(color.RGB8),
				})
			}
		}
	}

	return
}
func parsePlaceLayer(oData *LayoutData, place psd.Layer, memo *[]bool, H, W int) (elems []LayoutPlaceElement, err error) {

	// ドットの色毎の機能を取り出す
	img := place.Picker
	b := img.Bounds()

	test.Assert(b.Min.X >= 0)
	test.Assert(b.Min.Y >= 0)

	my := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		mx := 0
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)
			_, _, _, a := c.RGBA()
			// 透明は弾く
			if a > 0 {
				// 同色を塗りつぶしアルゴリズムで収集する
				elem, ok := collectByFloodFill(x, y, img, b, memo, mx, my, H, W)
				if ok {
					elems = append(elems, elem)
				}
			}

			mx++
		}

		my += W
	}

	return
}

var (
	// 他のドットを繋げる、表示はされない
	rgbConnector = color.RGB8{R: 255, G: 0, B: 255}
	// 今月分の日付群、1日、 2日、 ...
	rgbThisMonthDays = color.RGB8{R: 0, G: 255, B: 0}
	// 先月分の日付群、1日、 2日、 ...
	rgbPrevMonthDays = color.RGB8{R: 64, G: 255, B: 192}

	// ストロングトーン

	// 1月 - 12月
	rgbMonth1  = color.RGB8{R: 0, G: 149, B: 141}
	rgbMonth2  = color.RGB8{R: 0, G: 151, B: 219}
	rgbMonth3  = color.RGB8{R: 0, G: 98, B: 172}
	rgbMonth4  = color.RGB8{R: 27, G: 28, B: 128}
	rgbMonth5  = color.RGB8{R: 138, G: 1, B: 124}
	rgbMonth6  = color.RGB8{R: 214, G: 0, B: 119}
	rgbMonth7  = color.RGB8{R: 215, G: 0, B: 74}
	rgbMonth8  = color.RGB8{R: 215, G: 0, B: 15}
	rgbMonth9  = color.RGB8{R: 228, G: 142, B: 0}
	rgbMonth10 = color.RGB8{R: 243, G: 225, B: 0}
	rgbMonth11 = color.RGB8{R: 134, G: 184, B: 27}
	rgbMonth12 = color.RGB8{R: 0, G: 145, B: 64}
)

func collectByFloodFill(x, y int, img image.Image, b image.Rectangle, memo *[]bool, mx, my, H, W int) (elem LayoutPlaceElement, ok bool) {
	if (*memo)[mx+my] {
		return elem, false
	}

	c := img.At(x, y)
	rgb := color.RGB8Model.Convert(c).(color.RGB8)
	if rgb.Equal(rgbConnector) {
		return elem, false
	}

	rec(&elem, rgb, x, y, img, b, memo, mx, my, H, W)

	elem.Rgb = rgb
	return elem, true
}
func rec(elem *LayoutPlaceElement, parentRgb color.RGB8, x, y int, img image.Image, b image.Rectangle, memo *[]bool, mx, my, H, W int) {
	if (*memo)[mx+my] {
		return
	}

	c := img.At(x, y)
	rgb := color.RGB8Model.Convert(c).(color.RGB8)
	if rgb.Equal(parentRgb) {
		(*elem).XY = append((*elem).XY, point{X: int16(x), Y: int16(y)})
		(*memo)[mx+my] = true
	} else if rgb.Equal(rgbConnector) {
		// 通り道
		(*memo)[mx+my] = true
	} else {
		return
	}

	if x > b.Min.X {
		rec(elem, parentRgb, x-1, y, img, b, memo, mx-1, my, H, W)
	}
	if x < b.Max.X-1 {
		rec(elem, parentRgb, x+1, y, img, b, memo, mx+1, my, H, W)
	}
	if y > b.Min.Y {
		rec(elem, parentRgb, x, y-1, img, b, memo, mx, my-W, H, W)
	}
	if y < b.Max.Y-1 {
		rec(elem, parentRgb, x, y+1, img, b, memo, mx, my+W, H, W)
	}
}

// Save -- 配置図情報を書き出す
// NOTE: パース済みデータを読み書きして高速化
func (d *LayoutData) Save(w io.Writer) error {
	buf, err := msgpackEncodeLayout(d)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}
func msgpackEncodeLayout(d *LayoutData) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	_ = msgpack.NewEncoder(buf)

	// TODO: エンコード処理

	return buf, nil
}

// LoadPsdLayoutData -- 配置情報を読み込む
func LoadPsdLayoutData(r io.Reader) (data LayoutData, err error) {
	return msgpackDecodeLayout(r)
}
func msgpackDecodeLayout(r io.Reader) (data LayoutData, err error) {
	_ = msgpack.NewDecoder(r)

	// TODO: デコード処理

	return
}

// WriteSvg -- SVGとして書き出す
func (d LayoutData) WriteSvg(svgs svg.PixelaData, w io.Writer) {
	s := svgo.New(w)

	s.Startraw(fmt.Sprintf("viewBox=\"%d %d %d %d\"", d.MinX*10, d.MinY*10, (d.MaxX-d.MinX)*10, (d.MaxY-d.MinY)*10))
	// TODO: SVG書き出し処理
	for _, x := range d.Bg.Elems {
		s.Rect(int(x.X)*10, int(x.Y)*10, 9, 9, fmt.Sprintf("fill=\"%s\"", x.Rgb.ToColorCode()))
	}

	now := time.Now()
	thisYear := now.Year()
	thisMonth := int(now.Month())
	thisDay := 1

	prev := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
	prevYear := prev.Year()
	prevMonth := int(prev.Month())
	prevDay := 1
	for _, x := range d.Place.Elems {
		if x.Rgb.Equal(rgbThisMonthDays) {
			// 一致する日付を取り出す
			rgb, ok := func() (color.RGB8, bool) {
				for _, y := range svgs.Elems {
					if y.Date.Equal(thisYear, thisMonth, thisDay) {
						return y.Rgb, true
					}
				}

				return color.RGB8{}, false
			}()
			if ok {
				rect(s, x.XY, rgb)
			}

			thisDay++
		} else if x.Rgb.Equal(rgbPrevMonthDays) {
			rgb, ok := func() (color.RGB8, bool) {
				for _, y := range svgs.Elems {
					if y.Date.Equal(prevYear, prevMonth, prevDay) {
						return y.Rgb, true
					}
				}

				return color.RGB8{}, false
			}()
			if ok {
				rect(s, x.XY, rgb)
			}

			prevDay++
		} else if x.Rgb.Equal(rgbMonth1) {
			if thisMonth == 1 {
				rect(s, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth2) {
			if thisMonth == 2 {
				rect(s, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth3) {
			if thisMonth == 3 {
				rect(s, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth4) {
			if thisMonth == 4 {
				rect(s, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth5) {
			if thisMonth == 5 {
				rect(s, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth6) {
			if thisMonth == 6 {
				rect(s, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth7) {
			if thisMonth == 7 {
				rect(s, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth8) {
			if thisMonth == 8 {
				rect(s, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth9) {
			if thisMonth == 9 {
				rect(s, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth10) {
			if thisMonth == 10 {
				rect(s, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth11) {
			if thisMonth == 11 {
				rect(s, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth12) {
			if thisMonth == 12 {
				rect(s, x.XY, x.Rgb)
			}
		} else {
			log.Printf("unknown rgb: %s xy(len:%d [0]:%d %d)", x.Rgb.ToColorCode(), len(x.XY), x.XY[0].X, x.XY[0].Y)
		}
	}

	s.End()
}
func rect(s *svgo.SVG, xy []point, rgb color.RGB8) {
	for _, xy := range xy {
		s.Rect(int(xy.X)*10, int(xy.Y)*10, 9, 9, fmt.Sprintf("fill=\"%s\"", rgb.ToColorCode()))
	}
}
