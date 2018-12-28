package date

import (
	"time"
)

// Date -- 日付(年月日)クラス
type Date struct {
	time time.Time
}

const layout = "2006-01-02"

// GetString -- 日付を文字列で取得
func (dt Date) GetString() string {
	return dt.time.Format(layout)
}

// EqualYMD -- 年月日との一致チェック
func (dt Date) EqualYMD(year, month, day int) bool {
	return dt.time.Year() == year && int(dt.time.Month()) == month && dt.time.Day() == day
}

// Sub -- 日付との時差を返す
func (dt Date) Sub(o Date) time.Duration {
	return dt.time.Sub(o.time)
}

// From -- 日付をセット
func From(year, month, day int) (dt Date) {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

	return Date{time: t}
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
