package main

import (
	"fmt"
	"net/http"
	"os"

	dbUtil "slog-simple-blog/internal/database"
	logger "slog-simple-blog/internal/logger"
	server "slog-simple-blog/server"
)

func main() {
	sPort := os.Getenv("SLOG_SERVER_PORT")
	if sPort == "" {
		fmt.Println("No port specified, falling back to 8109")
		sPort = "8109"
	}
	app := &server.App{
		DB: dbUtil.InitPool(),
		Logger: logger.NewLogger(),
	}
	http.HandleFunc("/get_user_posts", app.HttpLogMiddleware(app.UserPostHanlder))
	http.HandleFunc("/new_user", app.HttpLogMiddleware(app.NewUserHandler))
	http.HandleFunc("/login", app.HttpLogMiddleware(app.LoginHandler))
	http.HandleFunc("/new_post", app.HttpLogMiddleware(app.AuthMiddleware((app.NewPostHandler))))
	http.HandleFunc("/logout", app.HttpLogMiddleware(app.AuthMiddleware(app.AuthManager.LogoutHandler)))
	go app.Logger.LogsWatchdog()
	if err := http.ListenAndServe(":" + sPort, nil); err != nil {
		panic(err)
	}
}
