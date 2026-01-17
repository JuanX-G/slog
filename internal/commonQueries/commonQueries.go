package commonQueries

import (
	"time"
)
type UserPostQueryRes struct {
	Content string `json:"content"`
	Title string `json:"title"`
	DatePosted time.Time `json:"date_posted"`
	Tags string `json:"tags"`
	Likes int32 `json:"likes"`
	ID int32 `json:"id"`
}

