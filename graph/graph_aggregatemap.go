package graph

import (
	"sort"

	"github.com/wordijp/pixelart/lib/map/slicemap"
)

// AggregateMap -- Pixela SVG集計データ
type AggregateMap = slicemap.MapStringInt

// Aggregate -- Pixela SVGの各色段階ごとの個数を集計する
func Aggregate(graph Data) AggregateMap {
	sort.SliceStable(graph.Elems, func(i, j int) bool {
		return graph.Elems[i].Count < graph.Elems[j].Count
	})

	// color毎のcount回数を集計する
	m := AggregateMap{}
	for _, x := range graph.Elems {
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
