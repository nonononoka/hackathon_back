package model

import "database/sql"

type User struct {
	ID        string         `json:"id"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	Bio       sql.NullString `json:"bio"`
	Image     sql.NullString `json:"image"`
	CreatedAt string         `json:"created_at"`
}
