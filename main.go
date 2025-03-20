package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/about", handleAbout)
	http.HandleFunc("/contact", handleContact)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
