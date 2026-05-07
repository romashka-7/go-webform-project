package main

import (
	"log"
	"net/http"

	"webform-go/internal/app"
	"webform-go/internal/config"
	apphttp "webform-go/internal/http"
	"webform-go/internal/http/handlers"
	"webform-go/internal/repository"
)

func main() {
	cfg := config.LoadConfig()

	db, err := app.NewBD(cfg)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewMySQLApplicationRepository(db)

	handlers.SetApplicationRepository(repo)

	router := apphttp.NewRouter()

	addr := ":" + cfg.ServerPort

	log.Println("Сервер запущен: http://localhost:" + addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
