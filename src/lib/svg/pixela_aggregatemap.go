package svg

import (
	"sort"

	"pixela_art/src/lib/map/slicemap"
)

// PixelaAggregateMap -- Pixela SVG集計データ
type PixelaAggregateMap = slicemap.MapStringInt

// AggregatePixelaData -- Pixela SVGの各色段階ごとの個数を集計する
func AggregatePixelaData(svgs PixelaData) PixelaAggregateMap {
	sort.SliceStable(svgs.Elems, func(i, j int) bool {
		return svgs.Elems[i].Count < svgs.Elems[j].Count
	})

	// color毎のcount回数を集計する
	m := PixelaAggregateMap{}
	for _, x := range svgs.Elems {
		color := x.Rgb.ToColorCode()

		idx, ok := m.Insert(color, 1)
		if !ok {
			_, val := m.NthRef(idx)
			*val++
		}
	}

	//for i, size := 0, m.Size(); i < size; i++ {
	//key, val := m.NthRef(i)
	//fmt.Printf("color(%s): count(%d)\n", key, *val)
	//}

	return m
}
