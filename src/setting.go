package main

import "time"

const location = "Asia/Tokyo"

func init() {
	// ロケーションの設定
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}
