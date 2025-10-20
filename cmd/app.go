package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"slog-simple-blog/internal/auth"
	dbUtil "slog-simple-blog/internal/database"
	pass "slog-simple-blog/internal/passwordUtils"
	passwordutils "slog-simple-blog/internal/passwordUtils"
)

type newPostQuery struct {
	Author string `json:"author"`
	Password string `json:"password"`
	Title string `json:"title"`
	Content string `json:"content"`
}

type userPostQuery struct {
	UserName string `json:"user_name"`
	Count int `json:"count"`
	Offset int `json:"offset"`
}

type userPostOffsetQuery struct { 
	UserName string `json:"user_name"`
	Count int `json:"count"`
 	Offset int `json:"offset"`

}

type userPostQueryRes struct {
	Content string `json:"content"`
	Title string `json:"title"`
	DatePosted time.Time `json:"date_posted"`
}

type newUserQuery struct {
	UserName string `json:"user_name"`
	Email string `json:"email"`
	DateCreated time.Time `json:"date_created"`
	Password string `json:"password"`
}

type newUserQueryRes struct {
	Response string `json:"response"`
}

var dbp *dbUtil.DB


func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithCancel(context.Background())
	user := r.FormValue("user")
	pass := r.FormValue("pass")

	dbRet, err := dbp.QueryForRow(ctx, "users", "user_name", user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var hashedPass string
	if len(dbRet) >= 5 {
		hashedPass = dbRet[4].(string)
	}
	if err := passwordutils.CheckPassword(hashedPass, pass); err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}

	auth.StoreLock.Lock()
	auth.SessionStore[token] = auth.Session{
		UserID:    user,
		ExpiresAt: time.Now().Add(auth.TokenTTL),
	}
	auth.StoreLock.Unlock()

	fmt.Fprintln(w, token)
}
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		auth.StoreLock.RLock()
		session, ok := auth.SessionStore[token]
		auth.StoreLock.RUnlock()

		if !ok || time.Now().After(session.ExpiresAt) {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
func userPostHanlder(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
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
	userVals, err := dbp.QueryForRow(ctx, "users", "user_name", reqB.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var authorId int32
	if len(userVals) > 0 {
		authorId = userVals[0].(int32)
	}
	posts, err := dbp.QueryCountOffset(ctx, reqB.Count, reqB.Offset, "posts", "author_id", authorId)
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


func newPostHandler(w http.ResponseWriter, req *http.Request) {
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
	userData, err := dbp.QueryForRow(ctx, "users", "user_name", reqB.Author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var userId int32
	var passHash string
	if len(userData) > 3 {
		userId = userData[0].(int32)
		passHash = userData[4].(string)
	}
	if err := pass.CheckPassword(passHash, reqB.Password); err != nil {
		http.Error(w, fmt.Sprint("password error", err.Error()),  http.StatusInternalServerError)
		return
	}
	var cols =  [...]string{"author_id", "content", "date_created", "tags"}

	err = dbp.InsertInto(ctx, "posts", cols[:], userId, reqB.Content, time.Now(), reqB.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Sucess"))
}


func newUserHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bd := req.Body
	bodyB, err := io.ReadAll(bd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var reqB newUserQuery
	err = json.Unmarshal(bodyB, &reqB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var res newUserQueryRes

	if len(reqB.UserName) < 2 {
		res.Response = "username too short, must be at least three characters"
		b, _ := json.Marshal(res)
		_, _ = w.Write(b)
		return
	}
	_, err = dbp.QueryForRow(ctx, "users", "email", reqB.Email)
	if err == nil {
		res.Response = "Account with that email exists"
		b, _ := json.Marshal(res)
		_, _ = w.Write(b)
		return
	} 

	_, err = dbp.QueryForRow(ctx, "users", "user_name", reqB.UserName)
	if err == nil {
		res.Response = "Account with that email exists"
		b, _ := json.Marshal(res)
		_, _ = w.Write(b)
		return
	}
	hash, err := pass.HashPassword(reqB.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var cols = [...]string{"user_name", "email", "date_created", "password"}
	err = dbp.InsertInto(ctx, "users", cols[:], reqB.UserName, reqB.Email, reqB.DateCreated, hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res.Response = "Account created succesfully"
	b, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}


func main() {
	dbp = dbUtil.InitPool()
	if err := dbp.ConfigureDB(); err != nil {
		panic(err)
	}
	http.HandleFunc("/get_user_posts", userPostHanlder)
	http.HandleFunc("/new_user", newUserHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/new_post", AuthMiddleware(newPostHandler))
	http.HandleFunc("/logout", AuthMiddleware(auth.LogoutHandler))
	http.ListenAndServe(":8080", nil)
}
