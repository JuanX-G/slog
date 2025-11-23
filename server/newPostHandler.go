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
	ctx, cancel := context.WithTimeout(req.Context(), time.Second * 5)
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
	if len(reqB.Title) < 3 {
		http.Error(w, "Title too short, must be at least 3 characters", http.StatusBadRequest)
		return
	}
	if len(reqB.Content) < 3 {
		http.Error(w, "Content too short, must be at least 3 characters", http.StatusBadRequest)
		return
	}
	if len(reqB.Content) > 800 {
		http.Error(w, "Content too long, must be at most 800 characters", http.StatusBadRequest)
		return
	}
	userData, err := a.DB.QueryForRow(ctx, "users", "user_name", reqB.Author)
	if err != nil {
		// todo fix error showing too much
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
