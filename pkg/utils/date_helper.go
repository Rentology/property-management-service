package utils

import (
	"fmt"
	"time"
)

func ParseDate(date *string) (string, error) {
	if date == nil {
		return "", fmt.Errorf("date is null")
	}
	newDate, err := time.Parse(time.RFC3339, *date)
	if err != nil {
		return "", err
	}
	formattedDate := newDate.Format("2006-01-02")
	return formattedDate, nil
}
