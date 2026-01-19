package federation

import (
	"net/http"
	"encoding/json"
	"io"
	"context"
	"time"

	queries "slog-simple-blog/internal/commonQueries"
)

type allPostsQuery struct {
	Since time.Time `json:"created_after"`
}

func (f Federation) FetchPosts(w http.ResponseWriter, req http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second * 5)
	defer cancel()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "error occured processing the request body", http.StatusBadRequest)
		return
	}
	var reqB allPostsQuery
	json.Unmarshal(body, &reqB)
	db := **f.db
	rows, err := db.SelectAll(ctx, "posts")
	if err != nil {
		http.Error(w, "error occured with the database", http.StatusBadRequest)
		return
	}
	cols := []string{"post_id"}
	var fRes []queries.UserPostQueryRes
	for _, v := range rows {
		var res queries.UserPostQueryRes
		res.Content = v[2].(string)
		res.DatePosted = v[3].(time.Time)
		res.Title = v[4].(string)
		res.Tags = v[5].(string)
		res.ID = v[0].(int32)
		if res.DatePosted.Before(reqB.Since) {continue}

		count, err := db.CountWhere(ctx, "post_likes", cols, res.ID)
		if err != nil {
			panic(err)
		}
		res.Likes = count
		fRes = append(fRes, res)
	}
	
	json.NewEncoder(w).Encode(fRes)
} 
