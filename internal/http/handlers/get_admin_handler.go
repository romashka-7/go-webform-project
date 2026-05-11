package handlers

import (
	"encoding/json"
	"net/http"
	"webform-go/internal/domain"
	"webform-go/internal/service"
)

func AdminGetApplications(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	svc := service.NewApplicationService(applicationRepo)

	applications, err := svc.GetAll()
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	response := domain.APIResponse{
		Status:  "success",
		Message: "Список заявок для администратора",
		Data:    applications,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
