package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const TOKEN_TTL = time.Hour
type AuthManager struct {
	sessionStore map[string]Session
	storeLock    sync.RWMutex
}

func(a *AuthManager) PurgeOutdatedTokens() {
	a.storeLock.RLock()
	var keysToPurge []string
	for k, v := range a.sessionStore {
		if time.Now().After(v.ExpiresAt) {
			keysToPurge = append(keysToPurge, k)
		}	
	}
	a.storeLock.RUnlock()
	a.storeLock.Lock()
	for _, e := range keysToPurge {
		delete(a.sessionStore, e)
	}
	a.storeLock.Unlock()
}

func NewAuthManager() *AuthManager {	
	return &AuthManager{
		sessionStore: make(map[string]Session),
	}
}

type Session struct {
	UserID    string
	ExpiresAt time.Time
}


func GenerateToken() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func(a *AuthManager) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Auth-Token")
	a.storeLock.Lock()
	_, f := a.sessionStore[token]
	if !f {
		fmt.Fprintln(w, "No such session found")
		return
	}
	delete(a.sessionStore, token)
	a.storeLock.Unlock()
	fmt.Fprintln(w, "Logged out")
}

func(a *AuthManager) Login(token, userID string) {
	a.storeLock.Lock()
	a.sessionStore[token] = Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(TOKEN_TTL),
	}
	a.storeLock.Unlock()
}

func(a *AuthManager) LookupToken(token string) (session Session, ok bool) {
	a.storeLock.RLock()
	session, ok = a.sessionStore[token]
	a.storeLock.RUnlock()
	return
}
