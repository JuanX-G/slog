package app

import (
	auth "slog-simple-blog/internal/auth"
	dbUtil "slog-simple-blog/internal/database"
	logger "slog-simple-blog/internal/logger"

)
type App struct  {
	DB *dbUtil.DB
	Logger *logger.Logger
	AuthManager *auth.AuthManager
}
