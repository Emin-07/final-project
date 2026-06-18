package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Emin-07/final-project/internal/pkg/hash"
)

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(os.Getenv("JWT_KEY")), nil
}

func Auth(next http.HandlerFunc, password string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		if len(password) > 0 {
			var tokenStr string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenStr = cookie.Value
			}
			valid := false

			claims := &jwt.RegisteredClaims{}
			_, err = jwt.ParseWithClaims(tokenStr, claims, keyFunc)
			if err == nil {
				if hash.CheckPasswordHash(password, claims.Subject) {
					valid = true
				}
			}

			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
