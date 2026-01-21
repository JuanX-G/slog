package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

)

type ForeignUserPostQueryRes struct {
	Content string `json:"content"`
	Title string `json:"title"`
	DatePosted time.Time `json:"date_posted"`
	Tags string `json:"tags"`
	OriginServerAddres string `json:"origin_server_addres"`
	ID int32 `json:"id"`
}

func(a *App) ForeignUserPostHanlder(w http.ResponseWriter, req *http.Request) {
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
	posts, err := a.DB.QueryCountOffset(ctx, reqB.Count, reqB.Offset, "foreign_posts", "origin_username", reqB.UserName)
	if err != nil {
		http.Error(w, "error occured in the database", http.StatusBadRequest)
		return
	}
	var fRes []ForeignUserPostQueryRes
	for _, v := range posts {
		var res ForeignUserPostQueryRes
		res.ID = v[0].(int32)
		res.Content = v[1].(string)
		res.DatePosted = v[2].(time.Time)
		res.Title = v[3].(string)
		res.Tags = v[4].(string)
		res.OriginServerAddres = v[5].(string)
		fRes = append(fRes, res)
	}

	json.NewEncoder(w).Encode(fRes)
}
