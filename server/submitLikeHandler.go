package app

import (
	"net/http"
	"context"
	"io"
	"encoding/json"
	"time"
)

type submitLikeQuery struct {
	PostID int32 `json:"post_id"`
	LikerID int32 `json:"liker_id"`
}

func(a App) SubmitLikeHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second * 5)
	defer cancel()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var reqB submitLikeQuery
	err = json.Unmarshal(body, &reqB)
	if err != nil {
		http.Error(w, "error reading the request content", http.StatusBadRequest)
		return
	}
	var cols = []string{"user_id", "post_id", "created_at"}
	err = a.DB.InsertInto(ctx, "post_likes", cols, reqB.LikerID, reqB.PostID, time.Now())	
	if err != nil {
		http.Error(w, "error executing the db query", http.StatusBadRequest)
		return
	}
}
