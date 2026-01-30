package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

// LocalTime 自定义时间类型，JSON 序列化为 "2006-01-02 15:04:05" 格式
type LocalTime time.Time

// MarshalJSON 序列化为 JSON
func (t LocalTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(TimeFormat))
	return []byte(stamp), nil
}

// UnmarshalJSON 从 JSON 反序列化
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	// 去掉引号
	str := string(data)[1 : len(data)-1]
	parsed, err := time.ParseInLocation(TimeFormat, str, time.Local)
	if err != nil {
		return err
	}
	*t = LocalTime(parsed)
	return nil
}

// Value 实现 driver.Valuer 接口，写入数据库
func (t LocalTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan 实现 sql.Scanner 接口，从数据库读取
func (t *LocalTime) Scan(v interface{}) error {
	if val, ok := v.(time.Time); ok {
		*t = LocalTime(val)
		return nil
	}
	return fmt.Errorf("cannot scan %T into LocalTime", v)
}

// Time 转换为 time.Time
func (t LocalTime) Time() time.Time {
	return time.Time(t)
}

// String 返回格式化字符串
func (t LocalTime) String() string {
	return time.Time(t).Format(TimeFormat)
}
