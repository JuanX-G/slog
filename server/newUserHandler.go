package app

import (
	"context"
	"net/http"
	"io"
	"encoding/json"
	"time"
	passwordUtil "slog-simple-blog/internal/passwordUtils"

)

type newUserQuery struct {
	UserName string `json:"user_name"`
	Email string `json:"email"`
	DateCreated time.Time 
	Password string `json:"password"`
}

type newUserQueryRes struct {
	Response string `json:"response"`
}

func(a *App) NewUserHandler(w http.ResponseWriter, req *http.Request) {
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
	reqB.DateCreated = time.Now()
	var res newUserQueryRes

	if len(reqB.UserName) < 2 {
		res.Response = "username too short, must be at least three characters"
		b, _ := json.Marshal(res)
		_, _ = w.Write(b)
		return
	}
	_, err = a.DB.QueryForRow(ctx, "users", "email", reqB.Email)
	if err == nil {
		res.Response = "Account with that email exists"
		b, _ := json.Marshal(res)
		_, _ = w.Write(b)
		return
	} 

	_, err = a.DB.QueryForRow(ctx, "users", "user_name", reqB.UserName)
	if err == nil {
		res.Response = "Account with that email exists"
		b, _ := json.Marshal(res)
		_, _ = w.Write(b)
		return
	}
	hash, err := passwordUtil.HashPassword(reqB.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var cols = [...]string{"user_name", "email", "date_created", "password"}
	err = a.DB.InsertInto(ctx, "users", cols[:], reqB.UserName, reqB.Email, reqB.DateCreated, hash)
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
