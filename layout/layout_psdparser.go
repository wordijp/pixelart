package layout

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"io"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	svgo "github.com/ajstarks/svgo"
	"github.com/oov/psd"
	"github.com/vmihailenco/msgpack"

	"github.com/wordijp/pixelart/graph"
	color "github.com/wordijp/pixelart/lib/color"
	test "github.com/wordijp/pixelart/lib/testify"

	"github.com/wordijp/pixelart/lib/date"
)

// Data -- 配置図用画像をパースする
type Data struct {
	Bg         DataBackgroundLayer
	Place      DataPlaceLayer
	MinX, MaxX int16
	MinY, MaxY int16
}

// DataBackgroundLayer -- 背景用レイヤーのパースデータ
type DataBackgroundLayer struct {
	Elems []DataBackgroundElement
}

// DataBackgroundElement -- 背景用レイヤーのドット情報
type DataBackgroundElement struct {
	X, Y int16
	Rgb  color.RGB8
}

// DataPlaceLayer -- 配置情報用レイヤーのパースデータ
type DataPlaceLayer struct {
	Elems []DataPlaceElement
}

// DataPlaceElement -- 配置情報用レイヤーの配置情報
type DataPlaceElement struct {
	XY  []point
	Rgb color.RGB8
}

type point struct {
	X, Y int16
}

// ParseLayoutPsd -- PSD画像をパースする
func ParseLayoutPsd(r io.Reader) (data Data, err error) {
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
func parseBackgroundLayer(bg psd.Layer) (data DataBackgroundLayer, err error) {
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
				data.Elems = append(data.Elems, DataBackgroundElement{
					X:   int16(x),
					Y:   int16(y),
					Rgb: color.RGB8Model.Convert(c).(color.RGB8),
				})
			}
		}
	}

	return
}
func parsePlaceLayer(oData *Data, place psd.Layer, memo *[]bool, H, W int) (elems []DataPlaceElement, err error) {

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

func collectByFloodFill(x, y int, img image.Image, b image.Rectangle, memo *[]bool, mx, my, H, W int) (elem DataPlaceElement, ok bool) {
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
func rec(elem *DataPlaceElement, parentRgb color.RGB8, x, y int, img image.Image, b image.Rectangle, memo *[]bool, mx, my, H, W int) {
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

// WriteLayoutData -- 配置図情報を書き出す
// NOTE: パース済みデータを読み書きして高速化
func (d *Data) WriteLayoutData(w io.Writer) error {
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
func msgpackEncodeLayout(d *Data) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	encoder := msgpack.NewEncoder(buf)

	// BackgroundLayer
	{
		length := uint32(len(d.Bg.Elems))

		err := encoder.EncodeUint32(length)
		if err != nil {
			return nil, err
		}

		// 各要素のビット情報をエンコード
		var bits uint64
		for i := uint32(0); i < length; i++ {
			x := uint64(d.Bg.Elems[i].X)
			y := uint64(d.Bg.Elems[i].Y)
			r := uint64(d.Bg.Elems[i].Rgb.R)
			g := uint64(d.Bg.Elems[i].Rgb.G)
			b := uint64(d.Bg.Elems[i].Rgb.B)

			bits = x<<48 | y<<32 | r<<16 | g<<8 | b

			if err := encoder.EncodeUint64(bits); err != nil {
				return nil, err
			}
		}
	}

	// PlaceLayer
	{
		length := uint32(len(d.Place.Elems))

		err := encoder.EncodeUint32(length)
		if err != nil {
			return nil, err
		}

		// 各要素のビット情報のエンコード
		for i := uint32(0); i < length; i++ {

			// xy
			{
				elemLength := uint32(len(d.Place.Elems[i].XY))

				err := encoder.EncodeUint32(elemLength)
				if err != nil {
					return nil, err
				}

				var bits uint32
				for j := uint32(0); j < elemLength; j++ {
					x := uint32(d.Place.Elems[i].XY[j].X)
					y := uint32(d.Place.Elems[i].XY[j].Y)

					bits = x<<16 | y

					if err := encoder.EncodeUint32(bits); err != nil {
						return nil, err
					}
				}
			}

			// Rgb
			{
				var bits uint32
				r := uint32(d.Place.Elems[i].Rgb.R)
				g := uint32(d.Place.Elems[i].Rgb.G)
				b := uint32(d.Place.Elems[i].Rgb.B)

				bits = r<<16 | g<<8 | b

				if err := encoder.EncodeUint32(bits); err != nil {
					return nil, err
				}
			}
		}
	}

	{
		var bits uint64
		minx := uint64(d.MinX)
		maxx := uint64(d.MaxX)
		miny := uint64(d.MinY)
		maxy := uint64(d.MaxY)

		bits = minx<<48 | maxx<<32 | miny<<16 | maxy

		if err := encoder.EncodeUint64(bits); err != nil {
			return nil, err
		}
	}

	return buf, nil
}

// LoadLayoutData -- 配置情報を読み込む
func LoadLayoutData(r io.Reader) (data Data, err error) {
	return msgpackDecode(r)
}
func msgpackDecode(r io.Reader) (data Data, err error) {
	decoder := msgpack.NewDecoder(r)

	// BackgroundLayer
	{
		var length uint32
		if length, err = decoder.DecodeUint32(); err != nil {
			return
		}

		data.Bg.Elems = make([]DataBackgroundElement, length, length)
		for i := uint32(0); i < length; i++ {
			// ビット情報からデコード
			var bits uint64
			if bits, err = decoder.DecodeUint64(); err != nil {
				return
			}
			data.Bg.Elems[i].X = int16(bits >> 48)
			data.Bg.Elems[i].Y = int16(bits >> 32)
			data.Bg.Elems[i].Rgb.R = uint8(bits >> 16)
			data.Bg.Elems[i].Rgb.G = uint8(bits >> 8)
			data.Bg.Elems[i].Rgb.B = uint8(bits >> 0)
		}
	}

	// PlaceLayer
	{
		var length uint32
		if length, err = decoder.DecodeUint32(); err != nil {
			return
		}

		data.Place.Elems = make([]DataPlaceElement, length, length)
		for i := uint32(0); i < length; i++ {

			// xy
			{
				var elemLength uint32
				if elemLength, err = decoder.DecodeUint32(); err != nil {
					return
				}

				data.Place.Elems[i].XY = make([]point, elemLength, elemLength)
				for j := uint32(0); j < elemLength; j++ {
					var bits uint32
					if bits, err = decoder.DecodeUint32(); err != nil {
						return
					}
					data.Place.Elems[i].XY[j].X = int16(bits >> 16)
					data.Place.Elems[i].XY[j].Y = int16(bits >> 0)
				}
			}

			// Rgb
			{
				var bits uint32
				if bits, err = decoder.DecodeUint32(); err != nil {
					return
				}
				data.Place.Elems[i].Rgb.R = uint8(bits >> 16)
				data.Place.Elems[i].Rgb.G = uint8(bits >> 8)
				data.Place.Elems[i].Rgb.B = uint8(bits >> 0)
			}
		}
	}

	{
		var bits uint64
		if bits, err = decoder.DecodeUint64(); err != nil {
			return
		}
		data.MinX = int16(bits >> 48)
		data.MaxX = int16(bits >> 32)
		data.MinY = int16(bits >> 16)
		data.MaxY = int16(bits >> 0)
	}

	return
}

// WriteSvgString -- SVGとして書き出す
func (d Data) WriteSvgString(svgs graph.Data, w io.Writer) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	s := svgo.New(bw)
	defer s.End()

	s.Startraw(fmt.Sprintf("viewBox=\"%d %d %d %d\"", d.MinX*10, d.MinY*10, (d.MaxX-d.MinX)*10, (d.MaxY-d.MinY)*10))
	{
		for _, x := range d.Bg.Elems {
			rect(bw, int(x.X)*10, int(x.Y)*10, 9, 9, x.Rgb)
		}
	}

	now := time.Now()
	thisYear := now.Year()
	thisMonth := int(now.Month())
	thisDay := 1

	prev := time.Date(thisYear, time.Month(thisMonth-1), 1, 0, 0, 0, 0, time.Local)
	prevYear := prev.Year()
	prevMonth := int(prev.Month())
	prevDay := 1
	for _, x := range d.Place.Elems {
		if x.Rgb.Equal(rgbThisMonthDays) {
			// 一致する日付を取り出す
			dt := date.From(thisYear, thisMonth, thisDay)
			idx := sort.Search(len(svgs.Elems), func(i int) bool {
				dur := svgs.Elems[i].Date.Sub(dt)
				return int(dur.Hours()) >= 0
			})

			if idx >= 0 && idx < len(svgs.Elems) && svgs.Elems[idx].Date.EqualYMD(thisYear, thisMonth, thisDay) {
				rects(bw, x.XY, svgs.Elems[idx].Rgb)
			}

			thisDay++
		} else if x.Rgb.Equal(rgbPrevMonthDays) {
			dt := date.From(prevYear, prevMonth, prevDay)
			idx := sort.Search(len(svgs.Elems), func(i int) bool {
				dur := svgs.Elems[i].Date.Sub(dt)
				return int(dur.Hours()) >= 0
			})

			if idx >= 0 && idx < len(svgs.Elems) && svgs.Elems[idx].Date.EqualYMD(prevYear, prevMonth, prevDay) {
				rects(bw, x.XY, svgs.Elems[idx].Rgb)
			}

			prevDay++
		} else if x.Rgb.Equal(rgbMonth1) {
			if thisMonth == 1 {
				rects(bw, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth2) {
			if thisMonth == 2 {
				rects(bw, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth3) {
			if thisMonth == 3 {
				rects(bw, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth4) {
			if thisMonth == 4 {
				rects(bw, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth5) {
			if thisMonth == 5 {
				rects(bw, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth6) {
			if thisMonth == 6 {
				rects(bw, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth7) {
			if thisMonth == 7 {
				rects(bw, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth8) {
			if thisMonth == 8 {
				rects(bw, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth9) {
			if thisMonth == 9 {
				rects(bw, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth10) {
			if thisMonth == 10 {
				rects(bw, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth11) {
			if thisMonth == 11 {
				rects(bw, x.XY, x.Rgb)
			}
		} else if x.Rgb.Equal(rgbMonth12) {
			if thisMonth == 12 {
				rects(bw, x.XY, x.Rgb)
			}
		} else {
			log.Printf("unknown rgb: %s xy(len:%d [0]:%d %d)", x.Rgb.ToColorCode(), len(x.XY), x.XY[0].X, x.XY[0].Y)
		}
	}
}
func rects(bw *bufio.Writer, xy []point, rgb color.RGB8) {
	for _, xy := range xy {
		rect(bw, int(xy.X)*10, int(xy.Y)*10, 9, 9, rgb)
	}
}
func rect(w *bufio.Writer, x, y, W, H int, rgb color.RGB8) {
	//fmt.Fprintf(w, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s" />`, int(x)*10, int(y)*10, 9, 9, rgb.ToColorCode())
	//fmt.Fprintln(w)
	w.WriteString(`<rect x="`)
	w.WriteString(strconv.Itoa(x))
	w.WriteString(`" y="`)
	w.WriteString(strconv.Itoa(y))
	w.WriteString(`" width="`)
	w.WriteString(strconv.Itoa(W))
	w.WriteString(`" height="`)
	w.WriteString(strconv.Itoa(H))
	w.WriteString(`" fill="`)
	w.WriteString(rgb.ToColorCode())
	w.WriteString(`" />`)
	w.WriteString("\n")
}
