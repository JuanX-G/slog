package app

import (
	"net/http"
	"context"
	"io"
	"encoding/json"
	"time"
)


func(a App) DeleteLikeHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second * 5)
	defer cancel()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var reqB LikeQuery
	err = json.Unmarshal(body, &reqB)
	if err != nil {
		http.Error(w, "error reading the request content", http.StatusBadRequest)
		return
	}
	var cols = []string{"user_id", "post_id"}
	err = a.DB.DeleteWhere(ctx, "post_likes", cols, reqB.LikerID, reqB.PostID)	
	if err != nil {
		http.Error(w, "error executing the db query", http.StatusBadRequest)
		return
	}
	w.Write([]byte("post succesfully unliked"))
}
