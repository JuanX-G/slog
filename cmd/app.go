package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"slog-simple-blog/internal/auth"
	dbUtil "slog-simple-blog/internal/database"
	logger "slog-simple-blog/internal/logger"
	server "slog-simple-blog/server"
)
f

func main() {
	sPort := os.Getenv("SLOG_SERVER_PORT")
	if sPort == "" {
		fmt.Println("No port specified, falling back to 8109")
		sPort = "8109"
	}
	app := &server.App{
		DB: dbUtil.InitPool(),
		Logger: logger.NewLogger(),
		AuthManager: auth.NewAuthManager(),
	}

	if err := app.DB.ConfigureDB(); err != nil {
		panic(err)
	}

	defer app.Logger.Close()
	go app.Logger.LogsWatchdog(true, false, "")
	go func () {
		ticker := time.NewTicker(time.Minute * 5)
		for range ticker.C {
			app.AuthManager.PurgeOutdatedTokens()
		}
	}()

	http.HandleFunc("/get_user_posts", app.HttpLogMiddleware(app.UserPostHanlder))
	http.HandleFunc("/new_user", app.HttpLogMiddleware(app.NewUserHandler))
	http.HandleFunc("/login", app.HttpLogMiddleware(app.LoginHandler))
	http.HandleFunc("/new_post", app.HttpLogMiddleware(app.AuthMiddleware((app.NewPostHandler))))
	http.HandleFunc("/logout", app.HttpLogMiddleware(app.AuthMiddleware(app.AuthManager.LogoutHandler)))
	http.HandleFunc("/get_user_description", app.HttpLogMiddleware(app.GetUserProfileHandler))
	http.HandleFunc("/submit_like", app.HttpLogMiddleware(app.AuthMiddleware(app.SubmitLikeHandler)))
	http.HandleFunc("/delete_like", app.HttpLogMiddleware(app.AuthMiddleware(app.DeleteLikeHandler)))
	if err := http.ListenAndServe(":" + sPort, nil); err != nil {
		panic(err)
	}
	app.Logger.LogString("Server running on: " + sPort + "\n")
}
