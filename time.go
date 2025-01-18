package ji

import (
	"fmt"
	"time"
)

// CNLoc 中国时区
// 项目加载初始化时可以使用 time.Local = ji.CNLoc
var CNLoc, _ = time.LoadLocation("Asia/Shanghai")

// DateTimeFormat 日期时间格式
const DateTimeFormat = "2006-01-02 15:04:05"

// DateFormat 日期格式
const DateFormat = "2006-01-02"

// TimeFormat 仅时间格式
const TimeFormat = "15:04:05"

// ParseTime 解析中国时间并返回
func ParseTime(str string) (time.Time, error) {
	// 使用 ParseInLocation 来解析时间字符串
	CNTime, err := time.ParseInLocation(DateTimeFormat, str, CNLoc)
	if err != nil {
		// 错误处理：返回一个带有详细信息的错误
		return time.Time{}, fmt.Errorf("时间解析失败: %v, 输入字符串: %s", err, str)
	}
	return CNTime, nil
}

// CurrentDateTime 当前日期时间
func CurrentDateTime() string {
	return time.Now().In(CNLoc).Format(DateTimeFormat)
}

// CurrentDate 当前日期
func CurrentDate() string {
	return time.Now().In(CNLoc).Format(DateFormat)
}

// CurrentTime 当前时间
func CurrentTime() string {
	return time.Now().In(CNLoc).Format(TimeFormat)
}

// DaysInMonth 返回指定时间所在月份的天数
func DaysInMonth(t time.Time) int {
	// 获取该月的第一天
	firstDayOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	// 获取下个月的第一天，并减去一天，得到当前月的最后一天
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)
	// 返回当前月的天数
	return lastDayOfMonth.Day()
}

// CurrentDaysInMonth 返回当前月份的天数
func CurrentDaysInMonth() int {
	// 获取当前时间并转换为中国时区
	return DaysInMonth(time.Now().In(CNLoc))
}

// CurrentTimestamp 返回当前时间的时间戳（秒级）
func CurrentTimestamp() int64 {
	return time.Now().In(CNLoc).Unix()
}

// TimestampToTime 将时间戳（秒级）转换为 time.Time
func TimestampToTime(ts int64) time.Time {
	return time.Unix(ts, 0).In(CNLoc)
}

// CurrentTimestampMilli 返回当前时间的时间戳（毫秒级）
func CurrentTimestampMilli() int64 {
	return time.Now().In(CNLoc).UnixNano() / int64(time.Millisecond)
}

// TimestampMilliToTime 将时间戳（毫秒级）转换为 time.Time
func TimestampMilliToTime(ts int64) time.Time {
	return time.Unix(0, ts*int64(time.Millisecond)).In(CNLoc)
}

// TimeToTimestamp 将 time.Time 转换为时间戳（秒级）
func TimeToTimestamp(t time.Time) int64 {
	return t.Unix()
}

// TimeToTimestampMilli 将 time.Time 转换为时间戳（毫秒级）
func TimeToTimestampMilli(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// IsValidDateTimeFormat 验证字符串是否是有效的日期时间格式
func IsValidDateTimeFormat(str string) bool {
	_, err := time.ParseInLocation(DateTimeFormat, str, CNLoc)
	return err == nil
}

// IsValidDateFormat 验证字符串是否是有效的日期格式
func IsValidDateFormat(str string) bool {
	_, err := time.ParseInLocation(DateFormat, str, CNLoc)
	return err == nil
}

// IsValidTimeFormat 验证字符串是否是有效的时间格式
func IsValidTimeFormat(str string) bool {
	_, err := time.ParseInLocation(TimeFormat, str, CNLoc)
	return err == nil
}

// FormatTimeToString 格式化时间为字符串
func FormatTimeToString(t time.Time, format string) string {
	return t.Format(format)
}

// ParseTimeWithCustomFormat 解析自定义格式的时间字符串
func ParseTimeWithCustomFormat(str string, format string) (time.Time, error) {
	return time.ParseInLocation(format, str, CNLoc)
}

// IsTimeAfter 验证 t1 是否在 t2 之后
func IsTimeAfter(t1, t2 time.Time) bool {
	return t1.After(t2)
}

// IsTimeBefore 验证 t1 是否在 t2 之前
func IsTimeBefore(t1, t2 time.Time) bool {
	return t1.Before(t2)
}

// TimestampDifference 计算两个时间戳之间的时间差
func TimestampDifference(ts1, ts2 int64) time.Duration {
	return time.Duration(ts2-ts1) * time.Second
}

// AddDurationToTime 向时间 t 添加指定的时长
func AddDurationToTime(t time.Time, duration time.Duration) time.Time {
	return t.Add(duration)
}

// SubtractDurationFromTime 从时间 t 减去指定的时长
func SubtractDurationFromTime(t time.Time, duration time.Duration) time.Time {
	return t.Add(-duration)
}

// ConvertTimeZone 将时间从一个时区转换到另一个时区
func ConvertTimeZone(t time.Time, fromZone, toZone *time.Location) time.Time {
	// 如果源时区和目标时区相同，直接返回原时间
	if fromZone == toZone {
		return t
	}
	// 将时间转换到目标时区
	return t.In(toZone)
}

// IsTimeEqual 验证两个时间是否相等（忽略时间的精度差异）
func IsTimeEqual(t1, t2 time.Time) bool {
	// 去掉秒以下的精度进行比较
	return t1.Truncate(time.Second).Equal(t2.Truncate(time.Second))
}

// GetDurationBetween 两个时间点之间的时长（返回小时、分钟和秒）
func GetDurationBetween(start, end time.Time) (int, int, int) {
	duration := end.Sub(start)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return hours, minutes, seconds
}
