package federation

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

)

type remotePosts struct {
	posts []ForeignUserPostQueryRes
	addres string
}
type ForeignUserPostQueryRes struct {
	Content string `json:"content"`
	Title string `json:"title"`
	DatePosted time.Time `json:"date_posted"`
	Tags string `json:"tags"`
	Likes int32 `json:"likes"`
	ID int32 `json:"id"`
	AuthorName string `json:"author_name"`
}
func(f Federation) FetchPostsFromAll() ([]remotePosts, error) {
	var allPosts []remotePosts
	for _, v := range f.knownServers {
		reqBody, _ := json.Marshal(map[string]time.Time{
			"created_after": v.lastSynced,
		})
		resp, err := http.Post(v.adress, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var posts []ForeignUserPostQueryRes
		if err = json.Unmarshal(body, &posts); err != nil {return nil, err}
		rPosts := remotePosts{posts: posts, addres: v.adress}
		allPosts = append(allPosts, rPosts)
	}
	return allPosts, nil
}

func(f Federation) SyncPostsFromAll(ctx context.Context) error {
	db := **f.db
	fPostsArr, err := f.FetchPostsFromAll()
	if err != nil {return err}
	for _, v := range fPostsArr {
		for _, e := range v.posts {
			err := db.InsertInto(ctx, "foreign_posts", []string{"content", "date_created", "title", "tags", "origin_server_addres", "origin_user_name"}, e.Content, e.DatePosted, e.Tags, e.Tags, v.addres, e.AuthorName)
			if err != nil {return err}
		}
	}
	return nil
}
