package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlerHome)
	http.HandleFunc("/about", handlerAbout)
	http.HandleFunc("/project", handlerProject)
	http.HandleFunc("/contact", handlerContact)

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
