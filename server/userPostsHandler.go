package app


import (
	"net/http"
	"context"
	"encoding/json"
	"io"
	"time"
)

type userPostQuery struct {
	UserName string `json:"user_name"`
	Count int `json:"count"`
	Offset int `json:"offset"`
}

type userPostQueryRes struct {
	Content string `json:"content"`
	Title string `json:"title"`
	DatePosted time.Time `json:"date_posted"`
}

func(a *App) UserPostHanlder(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second * 5)
	defer cancel()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var fRes []userPostQueryRes
	for _, v := range posts {
		var res userPostQueryRes
		res.Content = v[2].(string)
		res.DatePosted = v[3].(time.Time)
		res.Title = v[4].(string)
		fRes = append(fRes, res)
	}

	json.NewEncoder(w).Encode(fRes)
}
