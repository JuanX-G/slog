package app
import (
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"

	httpLog "slog-simple-blog/internal/httpLogging"
)
func(a *App) HttpLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := httpLog.HttpLogData {
			Method: r.Method,
			Referer: r.Header.Get("Referer"),
			UserAgent: r.Header.Get("User-Agent"),
			URL: r.URL.String(),
		}
		log.IP = RequestGetRemoteAddress(r)
		metrics := httpsnoop.CaptureMetrics(next, w, r)
		log.Code = metrics.Code
		log.Duration = metrics.Duration
		a.Logger.LogStruct(log)

	}
}

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func RequestGetRemoteAddress(r *http.Request) string {
		hdr := r.Header
		hdrRealIP := hdr.Get("X-Real-Ip")
		hdrForwardedFor := hdr.Get("X-Forwarded-For")
		if hdrRealIP == "" && hdrForwardedFor == "" {
			return ipAddrFromRemoteAddr(r.RemoteAddr)
		}
		
		return hdrRealIP
}
