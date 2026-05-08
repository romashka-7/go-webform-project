package handlers

import (
	"net/http"
	"webform-go/internal/config"
)

func CheckAdminAuth(w http.ResponseWriter, r *http.Request) bool {
	cfg := config.LoadConfig()

	login, password, ok := r.BasicAuth()
	if !ok || login != cfg.AdminLogin || password != cfg.AdminPassword {
		w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
		http.Error(w, "Требуется авторизация администратора", http.StatusUnauthorized)
		return false
	}

	return true
}
