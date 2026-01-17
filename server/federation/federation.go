package federation

import "slog-simple-blog/internal/database"

type Federation struct {
	knownServers map[string]federatedServer
	db **database.DB
}

type federatedServer struct {
	adress string
}
