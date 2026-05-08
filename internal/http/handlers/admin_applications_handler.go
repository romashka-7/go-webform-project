package handlers

import (
	"net/http"
)

func AdminApplicationsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		AdminGetApplications(w, r)
	case http.MethodDelete:
		AdminDeleteApplication(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

}
