package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	db, err := ConnectToDB() //name of the database
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// read file
	result, err := ReadDataFromFile("data/Data.dat")
	if err != nil {
		log.Fatal(err)
	}
	//Add data to database
	err = addDataToDB(db, result)
	if err != nil {
		log.Fatal(err)
	}

	// Use the result variable
	for _, data := range result {
		fmt.Printf("Energy: %.2f, Refractive Indicator: %.2f, Absorption Indicator: %.2f\n", data.Energy, data.RefractiveIndicator, data.AbsorptionIndicator)
	}

	// connect to database
	http.HandleFunc("/home", handleHome)
	http.HandleFunc("/about", handleAbout)
	http.HandleFunc("/contact", handleContact)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
