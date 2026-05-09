package handlers

import (
	"net/http"
)

func ApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		ApplicationAPIHandler(w, r)

	case http.MethodGet:
		if r.URL.Path == "/api/applications/" {
			GetApplicationsHandler(w, r)
			return
		}

		GetApplicationHandler(w, r)

	case http.MethodPut:
		UpdateApplicationHandler(w, r)

	case http.MethodDelete:
		DeleteApplicationHandler(w, r)

	default:
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
	}
}
