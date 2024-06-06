package model

type Tweet struct {
	ID       string `json:"id"`
	Body     string `json:"body"`
	PostedBy string `json:"postedby"`
}
