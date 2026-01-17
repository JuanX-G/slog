package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	queries "slog-simple-blog/internal/commonQueries"
)

type userPostQuery struct {
	UserName string `json:"user_name"`
	Count int `json:"count"`
	Offset int `json:"offset"`
}

type UserPostQueryRes struct {
	Content string `json:"content"`
	Title string `json:"title"`
	DatePosted time.Time `json:"date_posted"`
	Tags string `json:"tags"`
	Likes int32 `json:"likes"`
	ID int32 `json:"id"`
}

func(a *App) UserPostHanlder(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second * 5)
	defer cancel()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "error occured processing the request body", http.StatusBadRequest)
		return
	}
	var reqB userPostQuery
	json.Unmarshal(body, &reqB)
	if reqB.Count < 0 {
		reqB.Count = -1
	}
	userVals, err := a.DB.QueryForRow(ctx, "users", "user_name", reqB.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var authorId int32
	if len(userVals) > 0 {
		authorId = userVals[0].(int32)
	}
	posts, err := a.DB.QueryCountOffset(ctx, reqB.Count, reqB.Offset, "posts", "author_id", authorId)
	if err != nil {
		http.Error(w, "error occured in the database", http.StatusBadRequest)
		return
	}
	cols := []string{"post_id"}
	var fRes []queries.UserPostQueryRes
	for _, v := range posts {
		var res queries.UserPostQueryRes
		res.Content = v[2].(string)
		res.DatePosted = v[3].(time.Time)
		res.Title = v[4].(string)
		res.Tags = v[5].(string)
		res.ID = v[0].(int32)
		count, err := a.DB.CountWhere(ctx, "post_likes", cols, res.ID)
		if err != nil {
			panic(err)
		}
		res.Likes = count
		fRes = append(fRes, res)
	}

	json.NewEncoder(w).Encode(fRes)
}
