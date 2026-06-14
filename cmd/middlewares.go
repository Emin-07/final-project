package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		if len(app.config.password) > 0 {
			var tokenStr string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenStr = cookie.Value
			}
			valid := false

			claims := &jwt.RegisteredClaims{}
			_, err = jwt.ParseWithClaims(tokenStr, claims, app.keyFunc)
			if err == nil {
				if checkPasswordHash(app.config.password, claims.Subject) {
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
