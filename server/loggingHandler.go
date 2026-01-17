package app
import (
	"net/http"

	"github.com/felixge/httpsnoop"

	httpLog "slog-simple-blog/internal/httpLogging"
	"slog-simple-blog/internal/netHelpers"
)
func(a *App) HttpLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := httpLog.HttpLogData {
			Method: r.Method,
			Referer: r.Header.Get("Referer"),
			UserAgent: r.Header.Get("User-Agent"),
			URL: r.URL.String(),
		}
		log.IP = netHelpers.RequestGetRemoteAddress(r)
		metrics := httpsnoop.CaptureMetrics(next, w, r)
		log.Code = metrics.Code
		log.Duration = metrics.Duration
		a.Logger.LogStruct(log)

	}
}

