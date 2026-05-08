package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"webform-go/internal/domain"
	"webform-go/internal/service"
	"webform-go/internal/validation.go"
)

func UpdateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/applications/")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Некорректный ID заявки", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	var application domain.Application

	err = json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	err = validation.ValidateApplication(application)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	svc := service.NewApplicationService(applicationRepo)

	user, err := svc.GetUserBySessionID(cookie.Value)
	if err != nil {
		http.Error(w, "Сессия не найдена", http.StatusUnauthorized)
		return
	}

	if user.ApplicationID != id {
		http.Error(w, "Нельзя редактировать чужую заявку", http.StatusForbidden)
		return
	}

	updatedApplication, err := svc.Update(id, application)
	if err != nil {
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}

	response := domain.APIResponse{
		Status:  "success",
		Message: "Заявка успешно обновлена " + updatedApplication.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
