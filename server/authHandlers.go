package app

import (
	"net/http"
	"context"
	"fmt"
	"time"

	"slog-simple-blog/internal/auth"
	passwordutils "slog-simple-blog/internal/passwordUtils"
)
func(a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second * 5)
	defer cancel()
	user := r.FormValue("user")
	pass := r.FormValue("pass")
	dbRet, err := a.DB.QueryForRow(ctx, "users", "user_name", user)
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
	a.AuthManager.Login(token, user)
	fmt.Fprintln(w, token)
}
func(a *App) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		session, ok := a.AuthManager.LookupToken(token)
		if !ok || time.Now().After(session.ExpiresAt) {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
