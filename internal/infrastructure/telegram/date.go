package telegram

import (
	"fmt"
	"time"
)

const layout = "02.01.2006"

func parseRuDate(s string) (time.Time, error) {
	t, err := time.ParseInLocation(layout, s, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("некорректная дата: %v", err)
	}
	return t, nil
}
