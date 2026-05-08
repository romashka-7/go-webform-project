package handlers

import (
	"encoding/json"
	"net/http"

	"webform-go/internal/domain"
	"webform-go/internal/service"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var request domain.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	svc := service.NewApplicationService(applicationRepo)

	user, sessionID, err := svc.Login(request.Login, request.Password)
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	response := domain.APIResponse{
		Status:  "success",
		Message: "Авторизация успешна",
		Data: map[string]any{
			"user_id":        user.ID,
			"application_id": user.ApplicationID,
			"login":          user.Login,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
