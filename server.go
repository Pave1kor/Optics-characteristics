package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Настройки подключения к PostgreSQL
const (
	host     = "localhost"
	port     = 5432
	user     = "pavelkor"
	password = "1618"
	dbname   = "optics"
)

// Подключение к БД
func connectToDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}
	fmt.Println("Успешное подключение к базе данных")

	// Создание таблицы, если её нет
	query := `CREATE TABLE IF NOT EXISTS data (
		id SERIAL PRIMARY KEY,
		energy DOUBLE PRECISION, 
		reflection_coefficient DOUBLE PRECISION, 
		absorption_coefficient DOUBLE PRECISION
	)`
	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании таблицы: %w", err)
	}

	fmt.Println("Таблица проверена/создана")
	return db, nil
}

// Add data to baseData
func addDataToDB(db *sql.DB) error {
	//Load data
	result, err := readDataFromFile("data/Data.dat")
	if err != nil {
		return fmt.Errorf("ошибка чтения данных из файла: %w", err)
	}

	query := `INSERT INTO data (energy, reflection_coefficient, absorption_coefficient)
          VALUES ($1, $2, $3)`

	// Вставка данных
	for _, data := range result {
		_, err := db.Exec(query, data.Energy, data.ReflectionCoefficient, data.AbsorptionCoefficient)
		if err != nil {
			return fmt.Errorf("ошибка при вставке данных: %w", err)
		}
	}
	fmt.Println("Данные успешно добавлены")
	return nil
}

// Получение данных из БД
func getDataFromDB(db *sql.DB) ([]Data, error) {
	query := `SELECT energy, reflection_coefficient, absorption_coefficient FROM data`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе данных: %w", err)
	}
	defer rows.Close()

	var results []Data
	for rows.Next() {
		var data Data
		if err := rows.Scan(&data.Energy, &data.ReflectionCoefficient, &data.AbsorptionCoefficient); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании данных: %w", err)
		}
		results = append(results, data)
	}

	// Проверка на ошибки после итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после итерации по строкам: %w", err)
	}

	fmt.Println("Данные успешно загружены")
	return results, nil
}

func deleteDataFromDB(db *sql.DB) error {
	query := `DELETE FROM data`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка при удалении данных: %w", err)
	}
	fmt.Println("Данные успешно удалены")
	return nil
}
