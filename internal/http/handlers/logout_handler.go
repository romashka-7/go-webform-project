package handlers

import (
	"encoding/json"
	"net/http"

	"webform-go/internal/domain"
	"webform-go/internal/service"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	svc := service.NewApplicationService(applicationRepo)

	err = svc.Logout(cookie.Value)
	if err != nil {
		http.Error(w, "Ошибка выхода", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	response := domain.APIResponse{
		Status:  "success",
		Message: "Выход выполнен",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
