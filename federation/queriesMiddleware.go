package federation

import (
	"net/http"
	"time"
	"slog-simple-blog/internal/netHelpers"
)

func(f *Federation) HttpLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := netHelpers.RequestGetRemoteAddress(r)
		srv := federatedServer{adress: ip, lastSeen: time.Now()}
		f.knownServers[ip] = srv
		next(w, r)
	}
}
