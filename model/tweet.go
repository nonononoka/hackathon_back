package model

import "database/sql"

type Tweet struct {
	ID        string         `json:"id"`
	Body      string         `json:"body"`
	PostedBy  string         `json:"postedBy"`
	PostedAt  string         `json:"postedAt"`
	ReplyTo   sql.NullString `json:"replyTo"`
	LikeCount int            `json:"likeCount"`
	Tags      []string       `json:"tags"`
}
