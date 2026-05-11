package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"webform-go/internal/domain"
	"webform-go/internal/service"
	"webform-go/internal/validation"
)

func AdminUpdateApplication(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/admin/applications/")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Некорректный ID заявки", http.StatusBadRequest)
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

	updatedApplication, err := svc.Update(id, application)
	if err != nil {
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}

	response := domain.APIResponse{
		Status:  "success",
		Message: "Заявка успешно обновлена администратором",
		Data:    updatedApplication,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
