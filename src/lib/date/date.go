package date

import (
	"time"
)

// Date -- 日付扱うクラス
type Date struct {
	time time.Time
}

const layout = "2006-01-02"

// GetString -- 日付を文字列で取得
func (dt *Date) GetString() string {
	return dt.time.Format(layout)
}

// FromString -- 日付文字列をセット
func FromString(strDate string) (dt Date, err error) {
	t, err := time.Parse(layout, strDate)
	if err != nil {
		return
	}

	dt = Date{time: t}

	return
}
