package model

import "database/sql"

type Tweet struct {
	ID        string         `json:"id"`
	Body      string         `json:"body"`
	PostedBy  string         `json:"posted_by"`
	PostedAt  string         `json:"posted_at"`
	ReplyTo   sql.NullString `json:"reply_to"`
	LikeCount int            `json:"like_count"`
	Tags      []string       `json:"tags"`
}
