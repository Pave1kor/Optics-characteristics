package main

import (
	"database/sql"
	"fmt"
	"log"

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

type DataId struct {
	measurement_date   string
	measurement_number int
}
type DBManager struct {
	db *sql.DB
}

// Подключение к БД
func (manager *DBManager) connectToDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	manager.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	// Проверка соединения
	if err := manager.db.Ping(); err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %w", err)
	}
	fmt.Println("Успешное подключение к базе данных")
	return nil
}

// Get list of files
func (manager *DBManager) getListOfFiles() ([]DataId, error) {
	query := `SELECT measurement_date, measurement_number FROM files;`
	rows, err := manager.db.Query(query)
	dataSet := make([]DataId, 0)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе данных: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var data DataId
		if err := rows.Scan(&data.measurement_date, &data.measurement_number); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании данных: %w", err)
		}
		dataSet = append(dataSet, data)
	}
	return dataSet, nil
}

// Add data to baseData
func (manager *DBManager) addDataToDB(name string) error {
	// SQL-запрос для создания таблицы
	createTableQuery := `
CREATE TABLE IF NOT EXISTS measurements (
    id TEXT NOT NULL,
    type TEXT NOT NULL,
    measurement_date DATE NOT NULL,
    measurement_number INTEGER NOT NULL,
    column_name TEXT NOT NULL,
    x DOUBLE PRECISION NOT NULL,
    y DOUBLE PRECISION NOT NULL,
    PRIMARY KEY (id, column_name)
);`

	_, err := manager.db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}
	fmt.Println("Таблица успешно создана")
	//Load data - сделать универсальным
	result, title, err := readDataFromFile("data/Data.dat") //путь к файлу name.name
	if err != nil {
		return fmt.Errorf("ошибка чтения данных из файла: %w", err)
	}

	// Данные измерения, получить от пользователя
	measurementType := "temperature"
	measurementDate := "2025-03-28"

	// Получаем следующий номер измерения
	var nextNumber int
	query := `SELECT get_next_measurement_number($1, $2)`
	err = manager.db.QueryRow(query, measurementType, measurementDate).Scan(&nextNumber)
	if err != nil {
		log.Fatal(err)
	}

	// Генерируем ID
	id := fmt.Sprintf("%s_%s_%d", measurementType, measurementDate, nextNumber)

	// Вставляем каждую точку (x, y)
	insertQuery := `INSERT INTO measurements (id, type, measurement_date, measurement_number, column_name, x, y) 
					VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, point := range result {
		_, err := manager.db.Exec(insertQuery, id, measurementType, measurementDate, nextNumber, title, point.X, point.Y)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// Получение данных из БД
func (manager *DBManager) getDataFromDB(name string, title Title) ([]Data, error) {
	// добавить ключи
	query := fmt.Sprintf(`SELECT "%s", "%s" FROM %s`,
		title.X, title.Y, name)
	rows, err := manager.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе данных: %w", err)
	}
	defer rows.Close()

	var results []Data
	for rows.Next() {
		var data Data
		if err := rows.Scan(&data.X, &data.Y); err != nil {
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

// Удалениеданных из БД
func (manager *DBManager) deleteDataFromDB(name string) error {
	//добавить ключи
	query := fmt.Sprintf(`DELETE FROM %s`,
		name)
	_, err := manager.db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка при удалении данных: %w", err)
	}
	fmt.Println("Данные успешно удалены")
	return nil
}

// Удаление таблицы
func (manager *DBManager) dropTable(name string) error { //удалить по запросу
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, name)
	_, err := manager.db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка при удалении таблицы: %w", err)
	}
	fmt.Println("Таблица успешно удалена")
	return nil
}
