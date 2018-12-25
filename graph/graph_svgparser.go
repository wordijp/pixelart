// svgグラフデータをパースする

package graph

import (
	"io"
	"strconv"

	SVG "github.com/wordijp/svgparser"

	color "github.com/wordijp/pixelart/lib/color"
	date "github.com/wordijp/pixelart/lib/date"
)

// Data -- SVGパースデータ
type Data struct {
	Elems []DataElement
}

// DataElement -- 1要素
type DataElement struct {
	Date  date.Date  // svg: data-date
	Count int        // svg: data-count
	Rgb   color.RGB8 // svg: fill
}

// ParseCalendarGraphSvg -- Calendar Graph SVG文字列をパースする
func ParseCalendarGraphSvg(r io.Reader) (data Data, err error) {
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

func parseElement(oData *Data, svgElem *SVG.Element) error {
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

func parseRect(rectElem *SVG.Element) (elem DataElement, err error) {
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

	elem = DataElement{
		Date:  date,
		Count: count,
		Rgb:   rgb,
	}

	return
}
