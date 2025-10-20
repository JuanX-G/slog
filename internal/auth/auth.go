package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	UserID    string
	ExpiresAt time.Time
}

var (
	SessionStore = make(map[string]Session)
	StoreLock    sync.RWMutex
	TokenTTL     = time.Hour
)

func GenerateToken() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}



func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Auth-Token")
	StoreLock.Lock()
	delete(SessionStore, token)
	StoreLock.Unlock()
	fmt.Fprintln(w, "Logged out")
}
