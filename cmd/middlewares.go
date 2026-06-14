package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var tokenStr string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenStr = cookie.Value
				fmt.Println(tokenStr)
			}
			valid := false

			claims := &jwt.StandardClaims{}
			_, err = jwt.ParseWithClaims(tokenStr, claims, keyFunc)
			if err == nil {
				if checkPasswordHash(pass, claims.Subject) {
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
