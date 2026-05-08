package handlers

import (
	"encoding/json"
	"net/http"

	"webform-go/internal/domain"
	"webform-go/internal/service"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	svc := service.NewApplicationService(applicationRepo)

	user, err := svc.GetUserBySessionID(cookie.Value)
	if err != nil {
		http.Error(w, "Сессия не найдена", http.StatusUnauthorized)
		return
	}

	response := domain.APIResponse{
		Status:  "success",
		Message: "Пользователь авторизован",
		Data: map[string]any{
			"user_id":        user.ID,
			"application_id": user.ApplicationID,
			"login":          user.Login,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
