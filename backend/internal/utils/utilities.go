package utils

import (
	"database/sql"
	"strings"
)

// Helper: return a sql.NullString
func NullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: true, String: s}
}

func IsAllowedImageExtension(ext string) bool {
	ext = strings.ToLower(ext)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	return allowed[ext]
}

func DeleteUnusedImages() {

}
