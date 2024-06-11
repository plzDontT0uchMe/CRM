package utils

import (
	"database/sql"
	"strconv"
)

func ConvertStringToInt64(s string) int64 {
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
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
