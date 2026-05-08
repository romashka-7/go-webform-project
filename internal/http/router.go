package http

import (
	"net/http"

	"webform-go/internal/http/handlers"
)

func NewRouter() http.Handler {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/form", handlers.FormHandler)

	mux.HandleFunc("/api/applications", handlers.ApplicationsHandler)
	mux.HandleFunc("/api/applications/", handlers.ApplicationsHandler)

	mux.HandleFunc("/api/login", handlers.LoginHandler)

	mux.HandleFunc("/api/me", handlers.MeHandler)

	mux.HandleFunc("/api/logout", handlers.LogoutHandler)

	mux.HandleFunc("/admin/applications", handlers.AdminApplicationsHandler)
	mux.HandleFunc("/admin/applications/", handlers.AdminApplicationsHandler)
	return mux

}
