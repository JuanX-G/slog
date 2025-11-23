package httpLogger

import (
	"strconv"
	"strings"
	"time"
)
type HttpLogData struct {
	Method string 
	URL string
	Referer string
	UserAgent string
	IP string
	Code int
	Duration time.Duration
}

func (h HttpLogData) String() string {
	b := strings.Builder{}
	b.WriteString("HTTP request log\nURL: ")
	b.WriteString(h.URL + "\nMethod: ")
	b.WriteString(h.Method + "\nRemote address: ")
	b.WriteString(h.IP + "\nUser agent: ")
	b.WriteString(h.UserAgent + "\nReferer: ")
	b.WriteString(h.Referer + "\nCode: ")
	b.WriteString(strconv.Itoa(h.Code) + "\nDuration: ")
	b.WriteString(h.Duration.String() + "\n")

	return b.String()
}

