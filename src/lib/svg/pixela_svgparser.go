// pixelaのsvgデータをパースする

package svg

import (
	"io"
	"strconv"

	color "pixela_art/src/lib/color"
	date "pixela_art/src/lib/date"

	SVG "github.com/wordijp/svgparser"
)

// PixelaData -- SVGパースデータ
type PixelaData struct {
	Elems []PixelaDataElement
}

// PixelaDataElement -- 1要素
type PixelaDataElement struct {
	Date  date.Date  // svg: data-date
	Count int        // svg: data-count
	Rgb   color.RGB8 // svg: fill
}

// ParsePixelaSvg -- PixelaのSVG文字列をパースする
func ParsePixelaSvg(r io.Reader) (data PixelaData, err error) {
	svg, err := SVG.Parse(r, false)
	if err != nil {
		return
	}

	for _, child := range svg.Children {
		err = parseElement(&data, child)
		if err != nil {
			return
		}
	}

	return
}

func parseElement(oData *PixelaData, svgElem *SVG.Element) error {
	switch svgElem.Name {
	case "rect":
		parsed, err := parseRect(svgElem)
		if err == nil {
			oData.Elems = append(oData.Elems, parsed)
		}
	case "g":
		// groupは再帰的にパース
		for _, child := range svgElem.Children {
			err := parseElement(oData, child)
			if err != nil {
				return err
			}
		}
	default:
		// no-op
	}

	return nil
}

func parseRect(rectElem *SVG.Element) (elem PixelaDataElement, err error) {
	date, err := date.FromString(rectElem.Attributes["data-date"])
	if err != nil {
		return
	}
	count, err := strconv.Atoi(rectElem.Attributes["data-count"])
	if err != nil {
		return
	}
	rgb, err := color.RGB8FromColorCode(rectElem.Attributes["fill"])
	if err != nil {
		return
	}

	elem = PixelaDataElement{
		Date:  date,
		Count: count,
		Rgb:   rgb,
	}

	return
}
