package pojo

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

type JsonDate time.Time

func (t *JsonDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	str := string(data)
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse("2006-01-02", timeStr)
	*t = JsonDate(t1)
	return err
}

func (t JsonDate) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", time.Time(t).Format("2006-01-02"))
	return []byte(formatted), nil
}

func (t JsonDate) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format("2006-01-02"), nil
}

func (t *JsonDate) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = JsonDate(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *JsonDate) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}