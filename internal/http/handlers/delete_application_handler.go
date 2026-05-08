package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"webform-go/internal/domain"
	"webform-go/internal/service"
)

func DeleteApplicationHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/applications/")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Некорректный ID заявки", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	svc := service.NewApplicationService(applicationRepo)

	if err != nil {
		http.Error(w, "Сессия не найдена", http.StatusUnauthorized)
		return
	}

	err = svc.Delete(id)
	if err != nil {
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}

	response := domain.APIResponse{
		Status:  "success",
		Message: "Заявка успешно удалена",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
