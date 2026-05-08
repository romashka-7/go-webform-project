package middleware

import (
	"net/http"
	"os"
)

func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		login, password, ok := r.BasicAuth()

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)

			http.Error(
				w,
				"Требуется авторизация администратора",
				http.StatusUnauthorized,
			)

			return
		}

		adminLogin := os.Getenv("ADMIN_LOGIN")
		adminPassword := os.Getenv("ADMIN_PASSWORD")

		if login != adminLogin || password != adminPassword {

			w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)

			http.Error(
				w,
				"Неверный логин или пароль администратора",
				http.StatusUnauthorized,
			)

			return
		}

		next.ServeHTTP(w, r)
	})
}
