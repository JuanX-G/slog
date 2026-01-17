package federation

import (
	"net/http"
	server "slog-simple-blog/server/server"
)

func(f *Federation) HttpLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := server.RequestGetRemoteAddress(r)
		srv := federatedServer{adress: ip}
		f.knownServers[ip] = srv
		next(w, r)
	}
}
