package util

import "database/sql"

func StringPtr(s string) *string {
	return &s
}

func NullStringFrom(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}
