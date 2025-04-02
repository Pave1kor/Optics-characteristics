package repo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Pave1kor/Optics-characteristics/internal/app/models"
	read "github.com/Pave1kor/Optics-characteristics/internal/app/services"
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

type DBInterface interface {
	ConnectToDb() error
	GetListOfFiles() ([]models.DataId, error)
	AddDataToDB() error
	GetDataFromDB(title models.Title) ([]models.Data, error)
	DeleteDataFromDB() error
	DropTable() error
	Close() error
}

// Структура DBManager, которая реализует DBInterface
type DBManager struct {
	Db *sql.DB
}

// Конструктор для DBManager
func NewDBManager() *DBManager {
	return &DBManager{}
}

// Подключение к БД
func (manager *DBManager) ConnectToDb() error { //данные для подлючения внести в отдельный конфиг
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	manager.Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	// Проверка соединения
	if err := manager.Db.Ping(); err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %w", err)
	}
	fmt.Println("Успешное подключение к базе данных")
	return nil
}

// Get list of files
func (manager *DBManager) GetListOfFiles() ([]models.DataId, error) {
	query := `SELECT Date, Number FROM optics;`
	rows, err := manager.Db.Query(query)
	dataSet := make([]models.DataId, 0)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе данных: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var data models.DataId
		if err := rows.Scan(&data.Date, &data.Number); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании данных: %w", err)
		}
		dataSet = append(dataSet, data)
	}
	return dataSet, nil
}

// Add data to baseData
func (manager *DBManager) AddDataToDB() error {
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

	_, err := manager.Db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}
	fmt.Println("Таблица успешно создана")
	//Load data - сделать универсальным
	result, title, err := read.ReadDataFromFile("data/Data.dat") //путь к файлу name.name
	if err != nil {
		return fmt.Errorf("ошибка чтения данных из файла: %w", err)
	}
	// Данные измерения, получить от пользователя
	measurementType := "temperature"
	measurementDate := "2025-03-28"

	// Получаем следующий номер измерения
	var nextNumber int
	query := `SELECT get_next_measurement_number($1, $2)`
	err = manager.Db.QueryRow(query, measurementType, measurementDate).Scan(&nextNumber)
	if err != nil {
		log.Fatal(err)
	}

	// Генерируем ID
	id := fmt.Sprintf("%s_%s_%d", measurementType, measurementDate, nextNumber)

	// Вставляем каждую точку (x, y)
	insertQuery := `INSERT INTO measurements (id, type, measurement_date, measurement_number, column_name, x, y) 
					VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, point := range result {
		_, err := manager.Db.Exec(insertQuery, id, measurementType, measurementDate, nextNumber, title, point.X, point.Y)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// Получение данных из БД
func (manager *DBManager) GetDataFromDB(title models.Title) ([]models.Data, error) {
	// добавить ключи
	query := fmt.Sprintf(`SELECT "%s", "%s" FROM data`,
		title.X, title.Y)
	rows, err := manager.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе данных: %w", err)
	}
	defer rows.Close()

	var results []models.Data
	for rows.Next() {
		var data models.Data
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
func (manager *DBManager) DeleteDataFromDB() error {
	//добавить ключи
	// query := fmt.Sprintf(`DELETE FROM %s`,
	// 	name)
	query := `DELETE FROM data`
	_, err := manager.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка при удалении данных: %w", err)
	}
	fmt.Println("Данные успешно удалены")
	return nil
}

// Удаление таблицы
func (manager *DBManager) DropTable() error { //удалить по запросу
	query := fmt.Sprintf(`DROP TABLE IF EXISTS data`)
	_, err := manager.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка при удалении таблицы: %w", err)
	}
	fmt.Println("Таблица успешно удалена")
	return nil
}
func (manager *DBManager) Close() error {
	return manager.Db.Close()
}
