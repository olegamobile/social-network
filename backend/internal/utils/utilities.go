package utils

import "database/sql"

// Helper: return a sql.NullString
func NullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: true, String: s}
}
