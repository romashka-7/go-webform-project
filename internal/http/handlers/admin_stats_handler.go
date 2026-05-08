package handlers

import (
	"encoding/json"
	"net/http"

	"webform-go/internal/domain"
	"webform-go/internal/service"
)

func AdminStatsHandler(w http.ResponseWriter, r *http.Request) {
	if !CheckAdminAuth(w, r) {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	svc := service.NewApplicationService(applicationRepo)

	stats, err := svc.GetAdminStats()
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	response := domain.APIResponse{
		Status:  "success",
		Message: "Статистика администратора",
		Data:    stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
