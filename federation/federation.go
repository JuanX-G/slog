package federation

import (
	"slog-simple-blog/internal/database"
	"time"
)

type Federation struct {
	knownServers map[string]federatedServer
	db **database.DB
}

type federatedServer struct {
	adress string
	name string
	lastSeen time.Time
}
