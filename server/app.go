package app

import (
	auth "slog-simple-blog/internal/auth"
	appConfig "slog-simple-blog/internal/configUtil"
	dbUtil "slog-simple-blog/internal/database"
	logger "slog-simple-blog/internal/logger"
	federation "slog-simple-blog/federation"
)

type App struct  {
	DB dbUtil.Database
	Logger *logger.Logger
	AuthManager *auth.AuthManager
	Config *appConfig.AppConfig
	FederationMgr *federation.Federation
}
