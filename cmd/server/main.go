package main

import (
	"log"
	"net/http"

	"webform-go/internal/app"
	"webform-go/internal/config"
	apphttp "webform-go/internal/http"
	"webform-go/internal/http/handlers"
	"webform-go/internal/repository"
	"webform-go/internal/service"
)

func main() {
	cfg := config.LoadConfig() // load confing: port server, db, login/password admin

	db, err := app.NewBD(cfg) //connect to mysql

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // close db connection after main func

	repo := repository.NewMySQLApplicationRepository(db) // create repo which can work with db

	applicationService := service.NewApplicationService(repo) // create service layer which can work with repo

	handlers.SetApplicationRepository(repo) //transfer repo to handlers

	router := apphttp.NewRouter(applicationService) // create router

	addr := ":" + cfg.ServerPort // dynamic port server

	log.Println("Сервер запущен: http://localhost:" + addr) // info server start
	log.Fatal(http.ListenAndServe(addr, router))            // listen requests and serve them (brouser - client)
}
