package utils

import (
	"fmt"
	"time"
)

// Now 获取当前时间
func Now() time.Time {
	return time.Now()
}

// Since 计算时间差
func Since(t time.Time) time.Duration {
	return time.Since(t)
}

// FormatTime 格式化时间
func FormatTime(t time.Time, layout string) string {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return t.Format(layout)
}

// ParseTime 解析时间字符串
func ParseTime(timeStr string, layout string) (time.Time, error) {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return time.Parse(layout, timeStr)
}

// GetTimestamp 获取当前时间戳
func GetTimestamp() int64 {
	return time.Now().Unix()
}

// GetMillisecondTimestamp 获取当前毫秒时间戳
func GetMillisecondTimestamp() int64 {
	return time.Now().UnixMilli()
}

// GetDuration 获取格式化的时间差描述
func GetDuration(duration time.Duration) string {
	seconds := int(duration.Seconds())
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	switch {
	case days > 0:
		return fmt.Sprintf("%d天%d小时", days, hours%24)
	case hours > 0:
		return fmt.Sprintf("%d小时%d分钟", hours, minutes%60)
	case minutes > 0:
		return fmt.Sprintf("%d分钟%d秒", minutes, seconds%60)
	case seconds > 0:
		return fmt.Sprintf("%d秒", seconds)
	default:
		return fmt.Sprintf("%d毫秒", duration.Milliseconds())
	}
}