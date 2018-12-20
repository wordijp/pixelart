package image

import (
	"image/png"
	"io"

	"bytes"
	"encoding/gob"

	//"github.com/tinylib/msgp"

	color "pixela_art/src/lib/color"
)

// DotImageData -- ドット用画像をパースする
type DotImageData struct {
	Elems []DotImageElement
}

// DotImageElement -- ドット情報
type DotImageElement struct {
	X, Y int
	Rgb  color.RGB8
}

// ParseDotImage -- ドットアート用の画像をパースする
func ParseDotImage(r io.Reader) (data DotImageData, err error) {
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
				data.Elems = append(data.Elems, DotImageElement{
					X:   x,
					Y:   y,
					Rgb: color.RGB8Model.Convert(c).(color.RGB8),
				})
			}
		}
	}

	return
}

// Save -- ドット情報を書き出す
// NOTE: パース済みデータを読み書きして高速化
func (d *DotImageData) Save(w io.Writer) error {
	buf, err := gobEncode(d)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}
func gobEncode(d *DotImageData) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buf)

	err := gobEncodeDotImageData(d, encoder)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
func gobEncodeDotImageData(d *DotImageData, encoder *gob.Encoder) error {
	err := encoder.Encode(uint32(len(d.Elems)))
	if err != nil {
		return err
	}

	return encoder.Encode(d.Elems)
}
func gobEncodeDotImageElement(elem *DotImageElement, encoder *gob.Encoder) error {
	if err := encoder.Encode(elem); err != nil {
		return err
	}

	return nil
}

// LoadDotImageData -- ドット情報を読み込む
func LoadDotImageData(r io.Reader) (data DotImageData, err error) {
	return gobDecode(r)
}
func gobDecode(r io.Reader) (data DotImageData, err error) {
	decoder := gob.NewDecoder(r)

	var nelems uint32
	if err = decoder.Decode(&nelems); err != nil {
		return
	}

	data.Elems = make([]DotImageElement, nelems, nelems)
	err = decoder.Decode(&data.Elems)

	return
}
func gobDecodeDotImageElement(decoder *gob.Decoder) (elem DotImageElement, err error) {
	if err = decoder.Decode(&elem); err != nil {
		return
	}

	return
}
