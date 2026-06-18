package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func StringToken(passHashed string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		Subject:   passHashed,
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
