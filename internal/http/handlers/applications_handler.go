package handlers

import (
	"net/http"
)

func ApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		ApplicationAPIHandler(w, r)

	case http.MethodGet:
		GetApplicationsHandler(w, r)

	default:
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
	}
}
