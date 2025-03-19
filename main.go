package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// read file
	result, err := ReadDataFromFile("data/Data.dat")
	if err != nil {
		log.Fatal(err)
	}
	// Use the result variable
	for _, data := range result {
		fmt.Printf("Energy: %.2f, Refractive Indicator: %.2f, Absorption Indicator: %.2f\n", data.Energy, data.RefractiveIndicator, data.AbsorptionIndicator)
	}

	// connect to database
	db, err := ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
