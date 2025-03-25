package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/about", handleAbout)
	http.HandleFunc("/contact", handleContact)
	http.HandleFunc("/project", handleProject)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
