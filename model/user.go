package model

import "database/sql"

type User struct {
	ID        string         `json:"id"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	Bio       sql.NullString `json:"bio"`
	Image     sql.NullString `json:"image"`
	CreatedAt string         `json:"createdAt"`
	// userが今followしてるか
	IsFollowing bool `json:"isFollowing"`
	// userが今followされてるか
	IsFollowed bool `json:"isFollowed"`
}
