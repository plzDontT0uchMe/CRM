package utils

import (
	"CRM/go/usersService/internal/logger"
	"database/sql"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

func ConvertStringToInt64(s string) int64 {
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logger.CreateLog("error", err.Error())
		return 0
	}
	return result
}

func ConvertSQLNullStringToString(nullString sql.NullString) string {
	if nullString.Valid {
		return nullString.String
	}
	return ""
}

func ConvertTimestampToSQLNullTime(ts *timestamppb.Timestamp) sql.NullTime {
	if ts == nil {
		return sql.NullTime{
			Valid: false,
		}
	}
	return sql.NullTime{
		Time:  ts.AsTime(),
		Valid: true,
	}
}

func ConvertSQLNullTimeToTimestamp(nullTime sql.NullTime) *timestamppb.Timestamp {
	if nullTime.Valid {
		return timestamppb.New(nullTime.Time)
	}
	return nil
}

func ConvertStringToSQLNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
