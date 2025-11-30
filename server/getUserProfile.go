package app

import (
	"net/http"
	"context"
	"io"
	"encoding/json"
	"time"
)

type UserProfileQuery struct {
	UserName string `json:"user_name"`
}

type UserProfileQueryRes struct {
	Description string `json:"description"`
}

func (a *App) GetUserProfileHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), time.Second * 5)
	defer cancel()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var reqB userPostQuery
	json.Unmarshal(body, &reqB)

	userVals, err := a.DB.QueryForRow(ctx, "users", "user_name", reqB.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var userDescription string 
	if len(userVals) > 5 {
		userDescription = userVals[5].(string)
	}
	res := UserProfileQueryRes{Description: userDescription}
	json.NewEncoder(w).Encode(res)
} 
