package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Enter your database details
const (
	host     = "localhost"
	port     = 5432
	user     = "yourusername"
	password = "yourpassword"
	dbname   = "yourdbname"
)

// connect to database
func ConnectToDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// err = db.Ping()
	// if err != nil {
	// 	return nil, err
	// }
	fmt.Println("Successfully connected to the database!")
	return db, nil
}

// add data to database
func addDataToDB(db *sql.DB, result []Data) error {
	// create table
	query := `CREATE TABLE IF NOT EXISTS data (energy float, refractiveIndicator float, absorptionIndicator float)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	// insert data
	query = `INSERT INTO data (energy, refractiveIndicator, absorptionIndicator) VALUES `
	for _, data := range result {
		_, err := db.Exec(query, data.Energy, data.RefractiveIndicator, data.AbsorptionIndicator)
		if err != nil {
			return err
		}
	}
	return nil
}
