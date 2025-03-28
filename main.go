package main

import (
	"log"
	"net/http"

	"github.com/Pave1kor/Optics-characteristics/Internal/app/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/about", handlers.HandlerAbout)
	http.HandleFunc("/contact", handlers.HandlerContact)

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
