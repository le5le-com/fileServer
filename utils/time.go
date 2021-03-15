package utils

import (
	"time"
)

// TimeZh 中国习惯的时间格式
const TimeZh = "2006-01-02 15:04:05"

// GetGmtTimestamp 获取GMT时间戳
func GetGmtTimestamp() time.Time {
	timestamp := time.Now().UTC()
	gmt, _ := time.LoadLocation("GMT")
	gmtTime := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second(), 0, gmt)
	return time.Unix(gmtTime.Unix(), 0)
}
