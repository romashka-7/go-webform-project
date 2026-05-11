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

	adminAuth := middleware.AdminAuth
	fs := http.FileServer(http.Dir("web/static"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)

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
		adminAuth(
			http.HandlerFunc(handlers.AdminStatsHandler),
		),
	)
	mux.Handle(
		"/admin/panel",
		middleware.AdminAuth(
			http.HandlerFunc(handlers.AdminPanelHandler),
		),
	)

	mux.Handle(
		"/admin/applications",
		middleware.AdminAuth(
			http.HandlerFunc(handlers.AdminApplicationsHandler),
		),
	)

	mux.Handle(
		"/admin/applications/",
		middleware.AdminAuth(
			http.HandlerFunc(handlers.AdminApplicationsHandler),
		),
	)

	return mux

}
