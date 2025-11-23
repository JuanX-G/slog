package app

import (
	"net/http"
	"context"
	"encoding/json"
	"io"
	"time"
)

type newPostQuery struct {
	Author string `json:"author"`
	Title string `json:"title"`
	Content string `json:"content"`
}

func(a *App) NewPostHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bd := req.Body
	bodyB, err := io.ReadAll(bd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var reqB newPostQuery
	err = json.Unmarshal(bodyB, &reqB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userData, err := a.DB.QueryForRow(ctx, "users", "user_name", reqB.Author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var userId int32
	if len(userData) > 3 {
		userId = userData[0].(int32)
	}
	var cols =  [...]string{"author_id", "content", "date_created", "tags"}

	err = a.DB.InsertInto(ctx, "posts", cols[:], userId, reqB.Content, time.Now(), reqB.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Sucess"))
}
