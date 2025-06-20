package domain

import (
	"database/sql"
	"strings"
)

func nullToString(v sql.NullString) string {
	if v.Valid {
		return v.String
	}

	return ""
}

func stringToNull(s string) sql.NullString {
	if strings.TrimSpace(s) == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullToStringSlice(ns sql.NullString) []string {
	if !ns.Valid || strings.TrimSpace(ns.String) == "" {
		return []string{}
	}
	parts := strings.Split(ns.String, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func stringSliceToNull(labels []string) sql.NullString {
	if len(labels) == 0 {
		return sql.NullString{Valid: false}
	}
	cleaned := make([]string, 0, len(labels))
	for _, l := range labels {
		t := strings.TrimSpace(l)
		if t != "" {
			cleaned = append(cleaned, t)
		}
	}
	if len(cleaned) == 0 {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: strings.Join(cleaned, ","), Valid: true}
}
