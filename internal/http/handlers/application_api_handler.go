package handlers

import (
	"encoding/json"
	"net/http"
	"webform-go/internal/domain"
	"webform-go/internal/repository"
	"webform-go/internal/service"
	"webform-go/internal/validation.go"
)

func ApplicationAPIHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var application domain.Application

	err := json.NewDecoder(r.Body).Decode(&application)

	if err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	err = validation.ValidateApplication(application)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	repo := repository.NewMemoryApplicationRepository()

	svc := service.NewApplicationService(repo)

	createdApp, err := svc.Create(application)

	if err != nil {
		http.Error(w, "Error of Server", http.StatusInternalServerError)
		return
	}
	response := domain.APIResponse{
		Status:  "success",
		Message: "Заявка успешно принята: " + createdApp.Name,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
