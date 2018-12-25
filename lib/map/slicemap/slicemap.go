package slicemap

// sliceで表現するmap、要素数が少ない時にmapより高速
// C++のboost::flat_mapと同じ
// 通常のmapに対して以下の特徴がある
//   - keyが追加した順番となる
//   - 要素数が少ない場合、mapより操作が速い
//   - 要素数が多くなると遅くなる

// MapStringInt -- キーがstring、要素がintのslice_map
type MapStringInt struct {
	pairs []pair
}
type pair struct {
	key   string
	value int
}

// Insert -- key/valueを追加する
// @return 成功時) int: 追加時のidx, bool: true
//         失敗時) int: 重複した要素のidx, bool: false
func (m *MapStringInt) Insert(key string, value int) (int, bool) {
	if idx, hit := m.Find(key); hit {
		return idx, false
	}

	m.pairs = append(m.pairs, pair{key: key, value: value})

	return len(m.pairs) - 1, true
}

// Find -- keyを検索する
// @return int: 見つかった時のidx, bool: 見つかった時true
func (m *MapStringInt) Find(key string) (int, bool) {
	for i, kvp := range m.pairs {
		if kvp.key == key {
			return i, true
		}
	}

	return -1, false
}

// AtRef -- keyキーのvalueポインタを返す
// @return 見つからない時にnil
func (m *MapStringInt) AtRef(key string) *int {
	if idx, hit := m.Find(key); hit {
		return &m.pairs[idx].value
	}

	return nil
}

// NthRef -- idx番目のkey, value参照を返す
// @return 見つからない時にnil
func (m *MapStringInt) NthRef(idx int) (string, *int) {
	if idx < 0 {
		return "", nil
	}
	if idx >= len(m.pairs) {
		return "", nil
	}

	return m.pairs[idx].key, &m.pairs[idx].value
}

// Size -- key数を返す
func (m *MapStringInt) Size() int {
	return len(m.pairs)
}
