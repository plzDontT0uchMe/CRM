package utils

import (
	"database/sql"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"strings"
	"time"
)

func ConvertArrayStringToArraySQLNullString(array []string) []sql.NullString {
	var result []sql.NullString
	for _, item := range array {
		result = append(result, sql.NullString{String: item, Valid: true})
	}
	return result
}

func ConvertTimestampToNullString(timestamp *timestamppb.Timestamp) string {
	if timestamp.IsValid() {
		return time.Unix(timestamp.GetSeconds(), 0).Format("2006-01-02")
	}
	return ""
}

func ConvertStringToNullInt64(s string) int64 {
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return result
}

func ConvertTimestampToNullStringFull(timestamp *timestamppb.Timestamp) string {
	if timestamp.IsValid() {
		return time.Unix(timestamp.GetSeconds(), 0).Format("2006-01-02 15:04:05")
	}
	return ""
}

func ConvertStringToTimestamp(str string) (*timestamppb.Timestamp, error) {
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return nil, err
	}
	return timestamppb.New(t), nil
}

func ConvertStringToDate(str string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid format, please use YYYY-MM-DD")
	}
	return t, nil
}

func ConvertStringToTime(str string) (time.Time, error) {
	str = strings.Trim(str, "time.Date(")
	str = strings.TrimRight(str, ")")
	parts := strings.Split(str, ", ")

	if len(parts) != 8 {
		return time.Time{}, fmt.Errorf("invalid format")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, err
	}

	monthStr := strings.Trim(parts[1], "time.")
	month, err := time.Parse("Jan", monthStr)
	if err != nil {
		fmt.Printf("error parse month: %v", err)
		return time.Time{}, err
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil {
		fmt.Printf("error parse day: %v", err)
		return time.Time{}, err
	}

	hour, err := strconv.Atoi(parts[3])
	if err != nil {
		fmt.Printf("error parse hour: %v", err)
		return time.Time{}, err
	}

	minute, err := strconv.Atoi(parts[4])
	if err != nil {
		fmt.Printf("error parse minute: %v", err)
		return time.Time{}, err
	}

	second, err := strconv.Atoi(parts[5])
	if err != nil {
		fmt.Printf("error parse second: %v", err)
		return time.Time{}, err
	}

	nano, err := strconv.Atoi(parts[6])
	if err != nil {
		fmt.Printf("error parse nano: %v", err)
		return time.Time{}, err
	}

	return time.Date(year, month.Month(), day, hour, minute, second, nano, time.UTC), nil
}
