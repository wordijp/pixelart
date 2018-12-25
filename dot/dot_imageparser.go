package dot

import (
	"bytes"
	"fmt"
	"image/png"
	"io"

	svgo "github.com/ajstarks/svgo"
	"github.com/vmihailenco/msgpack"

	color "github.com/wordijp/pixelart/lib/color"
	"github.com/wordijp/pixelart/lib/math"

	"github.com/wordijp/pixelart/graph"
)

// Data -- ドットアート用データ
type Data struct {
	Elems      []DataElement
	MinX, MaxX int16
	MinY, MaxY int16
}

// DataElement -- ドット情報
type DataElement struct {
	X, Y int16
	Rgb  color.RGB8
}

// ParseDotPng -- ドットアート用の画像をパースする
func ParseDotPng(r io.Reader) (data Data, err error) {
	img, err := png.Decode(r)
	if err != nil {
		return
	}

	// ドット情報を取り出す
	b := img.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			c := img.At(x, y)
			_, _, _, a := c.RGBA()
			// 透明は弾く
			if a > 0 {
				data.Elems = append(data.Elems, DataElement{
					X:   int16(x),
					Y:   int16(y),
					Rgb: color.RGB8Model.Convert(c).(color.RGB8),
				})
			}
		}
	}

	data.MinX = int16(b.Min.X)
	data.MaxX = int16(b.Max.X)
	data.MinY = int16(b.Min.Y)
	data.MaxY = int16(b.Max.Y)

	return
}

// WriteDotData -- ドット情報を書き出す
// NOTE: パース済みデータを読み書きして高速化
func (d *Data) WriteDotData(w io.Writer) error {
	buf, err := msgpackEncode(d)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}
func msgpackEncode(d *Data) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	encoder := msgpack.NewEncoder(buf)

	length := uint32(len(d.Elems))

	err := encoder.EncodeUint32(length)
	if err != nil {
		return nil, err
	}

	// 各要素のビット情報をエンコード
	var bits uint64
	for i := uint32(0); i < length; i++ {
		x := uint64(d.Elems[i].X)
		y := uint64(d.Elems[i].Y)
		r := uint64(d.Elems[i].Rgb.R)
		g := uint64(d.Elems[i].Rgb.G)
		b := uint64(d.Elems[i].Rgb.B)

		bits = x<<48 | y<<32 | r<<16 | g<<8 | b

		if err := encoder.EncodeUint64(bits); err != nil {
			return nil, err
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

// LoadDotData -- ドット情報を読み込む
func LoadDotData(r io.Reader) (data Data, err error) {
	return msgpackDecode(r)
}

func msgpackDecode(r io.Reader) (data Data, err error) {
	decoder := msgpack.NewDecoder(r)

	var length uint32
	if length, err = decoder.DecodeUint32(); err != nil {
		return
	}

	data.Elems = make([]DataElement, length, length)
	for i := uint32(0); i < length; i++ {
		// ビット情報からデコード
		var bits uint64
		if bits, err = decoder.DecodeUint64(); err != nil {
			return
		}
		data.Elems[i].X = int16(bits >> 48)
		data.Elems[i].Y = int16(bits >> 32)
		data.Elems[i].Rgb.R = uint8(bits >> 16)
		data.Elems[i].Rgb.G = uint8(bits >> 8)
		data.Elems[i].Rgb.B = uint8(bits >> 0)
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

// Convert -- Graph SVGをもとにドットデータを変換する
func (d Data) Convert(g graph.Data) Data {
	agg := graph.Aggregate(g)

	cl := BuildColorLevels(d, agg)

	levelPercentage := []float64{0.1, 0.325, 0.55, 0.775, 1.0}
	return d.ApplyColorLevels(cl, levelPercentage)
}

// applyColorLevels -- ドットにカラーレベルを適用する
func (d Data) ApplyColorLevels(cl ColorLevels, colorLevelPercentage []float64) Data {
	dots := Data{
		Elems: make([]DataElement, len(d.Elems), len(d.Elems)),
		MinX:  d.MinX,
		MaxX:  d.MaxX,
		MinY:  d.MinY,
		MaxY:  d.MaxY,
	}

	for i, length := 0, len(d.Elems); i < length; i++ {
		// 色相固定で、彩度・明度に割合適用

		// HSVで計算し
		hsv := d.Elems[i].Rgb.ToHSV()
		// 彩度は0へ
		hsv.S = uint8(math.Lerpf(float64(hsv.S), 0.0, 1.0-colorLevelPercentage[cl.levels[i]]))
		// 明度は100へ
		hsv.V = uint8(math.Lerpf(float64(hsv.V), 100.0, 1.0-colorLevelPercentage[cl.levels[i]]))

		// RGBに戻す
		dots.Elems[i].Rgb = hsv.ToRGB8()

		dots.Elems[i].X = d.Elems[i].X
		dots.Elems[i].Y = d.Elems[i].Y
	}

	return dots
}

// WriteSvgString -- SVGとして書き出す
func (d Data) WriteSvgString(w io.Writer) {
	s := svgo.New(w)

	s.Startraw(fmt.Sprintf("viewBox=\"%d %d %d %d\"", d.MinX*10, d.MinY*10, (d.MaxX-d.MinX)*10, (d.MaxY-d.MinY)*10))
	for _, x := range d.Elems {
		s.Rect(int(x.X)*10, int(x.Y)*10, 9, 9, fmt.Sprintf("fill=\"%s\"", x.Rgb.ToColorCode()))
	}
	s.End()
}
