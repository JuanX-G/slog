package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	federation "slog-simple-blog/federation"
	auth "slog-simple-blog/internal/auth"
	config "slog-simple-blog/internal/configUtil"
	dbUtil "slog-simple-blog/internal/database"
	logger "slog-simple-blog/internal/logger"
	server "slog-simple-blog/server"
)

func main() {
	logger, err := logger.NewLogger(false, "")
	if err != nil {panic(err)}
	app := &server.App{
		Config: config.NewAppConfig(),
		DB: &dbUtil.DB{},
		Logger: logger,
		AuthManager: auth.NewAuthManager(),
	}
	dbp := &app.DB
	app.FederationMgr = federation.NewFederationMgr(&dbp)
	app.Config.ParseConfig("./config.txt")
	sPort := app.Config.ServerPort
	app.DB, err = dbUtil.InitPool(app.Config.DBport)
	if err != nil {
		panic(err)
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
	
	var remoteSyncCancel func()
	defer remoteSyncCancel()
	go func() {
		secsInInt64 := int64(app.Config.ForeignSyncIntervalSecs)
		ticker := time.NewTicker(time.Duration(secsInInt64 * 1000000000))
		for range ticker.C {
			var ctx context.Context
			ctx, remoteSyncCancel = context.WithTimeout(context.Background(), time.Minute * 10)
			app.FederationMgr.SyncPostsFromAll(ctx)
		}
	}()

	if sPort == "" {
		fmt.Println("No port specified, falling back to 8109")
		sPort = "8109"
	}
	http.HandleFunc("/get_user_posts", app.HttpLogMiddleware(app.UserPostHanlder))
	http.HandleFunc("/new_user", app.HttpLogMiddleware(app.NewUserHandler))
	http.HandleFunc("/login", app.HttpLogMiddleware(app.LoginHandler))
	http.HandleFunc("/new_post", app.HttpLogMiddleware(app.AuthMiddleware((app.NewPostHandler))))
	http.HandleFunc("/logout", app.HttpLogMiddleware(app.AuthMiddleware(app.AuthManager.LogoutHandler)))
	http.HandleFunc("/get_user_description", app.HttpLogMiddleware(app.GetUserProfileHandler))
	http.HandleFunc("/submit_like", app.HttpLogMiddleware(app.AuthMiddleware(app.SubmitLikeHandler)))
	http.HandleFunc("/delete_like", app.HttpLogMiddleware(app.AuthMiddleware(app.DeleteLikeHandler)))

	http.HandleFunc("/federation/get_since", app.HttpLogMiddleware(app.FederationMgr.HttpLogMiddleware(app.DeleteLikeHandler)))
	app.Logger.LogString("Server running on: " + sPort + "\n")
	if err := http.ListenAndServe(":" + sPort, nil); err != nil {
		panic(err)
	}
}
