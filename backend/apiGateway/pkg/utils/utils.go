package utils

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

func ConvertInt64ToSQLNullString(val int64) string {
	if val == 0 {
		return ""
	}
	return strconv.FormatInt(val, 10)
}

func IndexOfStrings(slice []string, element string) int {
	for i, v := range slice {
		if v == element {
			return i
		}
	}
	return -1
}

func Contains(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
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
