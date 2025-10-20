package passwordutils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {
	if len(pass) > 72 {
		return "", errors.New("Passowrd to long, max length is 71")
	} 
	b, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func CheckPassword(hash, pass string) error { 
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}
