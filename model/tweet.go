package model

import "database/sql"

type Tweet struct {
	ID            string         `json:"id":"id"`
	Body          string         `json:"body" :"body"`
	PostedBy      string         `json:"postedBy" :"posted_by"`
	PostedByImage sql.NullString `json:"postedByImage" :"posted_by_image"`
	PostedByName  string         `json:"postedByName" :"posted_by_name"`
	PostedAt      string         `json:"postedAt" :"posted_at"`
	ReplyTo       sql.NullString `json:"replyTo" :"reply_to"`
	LikeCount     int            `json:"likeCount" :"like_count"`
	ReplyCount    int            `json:"replyCount" :"reply_count"`
	Tags          []string       `json:"tags" :"tags"`
	IsFaved       bool           `json:"isFaved" :"is_faved"`
}
