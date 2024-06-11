package utils

import (
	"database/sql"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertSQLNullTimeToTimestamp(nullTime sql.NullTime) *timestamppb.Timestamp {
	if nullTime.Valid {
		return timestamppb.New(nullTime.Time)
	}
	return nil
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

func ConvertInt64ToSQLNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}
