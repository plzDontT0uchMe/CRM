package utils

import "database/sql"

func NullStringToString(nullString sql.NullString) string {
	if nullString.Valid {
		return nullString.String
	}
	return ""
}
