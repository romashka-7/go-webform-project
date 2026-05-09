package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"webform-go/internal/domain"
	"webform-go/internal/service"
)

func GetApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	svc := service.NewApplicationService(applicationRepo)

	application, err := svc.GetAll()

	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка получения вызова", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status string               `json:"status"`
		Data   []domain.Application `json:"data"`
	}{
		Status: "success",
		Data:   application,
	}

	json.NewEncoder(w).Encode(response)
}
