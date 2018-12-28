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

// EqualYMD -- 年月日との一致チェック
func (dt Date) EqualYMD(year, month, day int) bool {
	return dt.time.Year() == year && int(dt.time.Month()) == month && dt.time.Day() == day
}

// FromString -- 日付文字列をセット
func FromString(strDate string) (dt Date, err error) {
	t, err := time.ParseInLocation(layout, strDate, time.Local)
	if err != nil {
		return
	}

	dt = Date{time: t}
	return
}
