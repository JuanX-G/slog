package netHelpers

import (
	"net/http"
	"strings"
)

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
