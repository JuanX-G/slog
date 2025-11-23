package app

import (
	dbUtil "slog-simple-blog/internal/database"
	logger "slog-simple-blog/internal/logger"
	auth "slog-simple-blog/internal/auth"
)
type App struct  {
	DB *dbUtil.DB
	Logger *logger.Logger
	AuthManager auth.AuthManager
}
