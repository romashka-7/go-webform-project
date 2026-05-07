package handlers

import (
	"encoding/json"
	"net/http"
	"webform-go/internal/domain"
	"webform-go/internal/repository"
	"webform-go/internal/service"
	"webform-go/internal/validation.go"
)

var applicationRepo repository.ApplicationRepository

func SetApplicationRepository(repo repository.ApplicationRepository) {
	applicationRepo = repo
}

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

	svc := service.NewApplicationService(applicationRepo)

	createdApp, login, password, err := svc.Create(application)

	if err != nil {
		http.Error(w, "Error of Server", http.StatusInternalServerError)
		return
	}
	response := domain.APIResponse{
		Status:  "success",
		Message: "Заявка успешно принята: " + createdApp.Name,
		Data: map[string]string{
			"login":    login,
			"password": password,
		},
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
