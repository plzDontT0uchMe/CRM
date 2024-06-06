package utils

import (
	"database/sql"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
