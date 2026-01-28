package entity

import "database/sql"

type Category struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}
