package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"webform-go/internal/domain"
	"webform-go/internal/service"
)

type contextKey string

const userContextKey contextKey = "user"

func Auth(applicationService *service.ApplicationService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
				return
			}

			user, err := applicationService.GetUserBySessionID(cookie.Value)
			if err != nil {
				http.Error(w, "Недействительная сессия", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r)
		if !ok {
			http.Error(w, "Пользователь не найден в контексте", http.StatusUnauthorized)
			return
		}

		idPart := strings.TrimPrefix(r.URL.Path, "/api/applications/")
		applicationID, err := strconv.Atoi(idPart)
		if err != nil {
			http.Error(w, "Некорректный ID заявки", http.StatusBadRequest)
			return
		}

		if user.ApplicationID != applicationID {
			http.Error(w, "Нет доступа к чужой заявке", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UserFromContext(r *http.Request) (domain.User, bool) {
	user, ok := r.Context().Value(userContextKey).(domain.User)
	return user, ok
}
