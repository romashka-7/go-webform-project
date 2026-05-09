package handlers

import (
	"encoding/json"
	"net/http"

	"webform-go/internal/domain"
	"webform-go/internal/service"
)

func writeJSON(w http.ResponseWriter, statusCode int, response domain.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, domain.APIResponse{
			Status:  "error",
			Message: "Метод не поддерживается",
		})
		return
	}

	var request domain.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, domain.APIResponse{
			Status:  "error",
			Message: "Невалидный JSON",
		})
		return
	}

	svc := service.NewApplicationService(applicationRepo)

	user, sessionID, err := svc.Login(request.Login, request.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, domain.APIResponse{
			Status:  "error",
			Message: "Неверный логин или пароль",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusOK, domain.APIResponse{
		Status:  "success",
		Message: "Авторизация успешна",
		Data: map[string]any{
			"user_id":        user.ID,
			"application_id": user.ApplicationID,
			"login":          user.Login,
		},
	})
}
