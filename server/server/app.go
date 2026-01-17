package app

import (
	auth "slog-simple-blog/internal/auth"
	dbUtil "slog-simple-blog/internal/database"
	logger "slog-simple-blog/internal/logger"
	appConfig "slog-simple-blog/internal/configUtil"

)
type App struct  {
	DB dbUtil.Database
	Logger *logger.Logger
	AuthManager *auth.AuthManager
	Config *appConfig.AppConfig
}
