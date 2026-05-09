package http

import (
	"net/http"

	"webform-go/internal/http/handlers"
	"webform-go/internal/http/middleware"
	"webform-go/internal/service"
)

func NewRouter(applicationService *service.ApplicationService) http.Handler {
	mux := http.NewServeMux()

	authMiddleware := middleware.Auth(applicationService)

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/form", handlers.FormHandler)

	mux.HandleFunc("/api/login", handlers.LoginHandler)

	mux.HandleFunc("/api/me", handlers.MeHandler)

	mux.HandleFunc("/api/logout", handlers.LogoutHandler)

	mux.Handle(
		"/api/applications/",
		authMiddleware(
			middleware.RequireOwner(
				http.HandlerFunc(handlers.ApplicationsHandler),
			),
		),
	)

	mux.Handle(
		"/admin/stats",
		middleware.AdminAuth(
			http.HandlerFunc(handlers.AdminStatsHandler),
		),
	)
	return mux

}
