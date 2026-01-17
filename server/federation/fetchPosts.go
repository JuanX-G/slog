package federation

import (
	"net/http"
	"encoding/json"
	"io"
	"context"
	"time"

	server "slog-simple-blog/server/server"
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
	rows, err := db.SelectAllWhere(ctx, "posts", []string{"date_created"}, reqB.Since)
	if err != nil {
		http.Error(w, "error occured with the database", http.StatusBadRequest)
		return
	}
	cols := []string{"post_id"}
	var fRes []server.UserPostQueryRes
	for _, v := range rows {
		var res server.UserPostQueryRes
		res.Content = v[2].(string)
		res.DatePosted = v[3].(time.Time)
		res.Title = v[4].(string)
		res.Tags = v[5].(string)
		res.ID = v[0].(int32)
		count, err := db.CountWhere(ctx, "post_likes", cols, res.ID)
		if err != nil {
			panic(err)
		}
		res.Likes = count
		fRes = append(fRes, res)
	}
	
	json.NewEncoder(w).Encode(fRes)
} 
