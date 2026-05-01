package main

import (
	"log"
	"net/http"

	apphttp "webform-go/internal/http"
)

func main() {
	
	router := apphttp.NewRouter()

	log.Println("Сервер запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}