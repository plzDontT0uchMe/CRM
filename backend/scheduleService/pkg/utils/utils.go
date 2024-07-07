package utils

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertTimestampToString(timestamp *timestamppb.Timestamp) string {
	return timestamp.AsTime().Format("2006-01-02")
}
