package repo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Pave1kor/Optics-characteristics/internal/app/models"
	read "github.com/Pave1kor/Optics-characteristics/internal/app/services"
	config "github.com/Pave1kor/Optics-characteristics/internal/config"
	_ "github.com/lib/pq"
)

type DBInterface interface {
	ConnectToDb() error
	GetListOfTables() ([]models.TableName, error)
	AddDataToDB() error
	GetDataFromDB(title models.Title) ([]models.Data, error)
	DeleteDataFromDB() error
	DropTable() error
	GenerateMeasurementID(measurementType, measurementDate string) (string, int, error)
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
	cfg := config.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "pavelkor",
		Password: "1618",
		DBName:   "optics",
		SSLMode:  "disable",
	}
	var err error
	conn := cfg.ConnString()
	manager.Db, err = sql.Open("postgres", conn)
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
func (manager *DBManager) GetListOfTables() ([]models.TableName, error) {
	rows, err := manager.Db.Query(models.GetListOfTablesQuery)
	list := []models.TableName{}
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	// Вывод списка таблиц
	fmt.Println("Список таблиц:")
	for rows.Next() {
		var tableName models.TableName
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		list = append(list, tableName)
	}
	return list, nil
}

// Add data to baseData
func (manager *DBManager) AddDataToDB() error {
	// SQL-запрос для создания таблицы
	_, err := manager.Db.Exec(models.CreateTableQuery)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}
	fmt.Println("Таблица успешно создана")
	//Load data - сделать универсальным
	result, title, err := read.ReadDataFromFile("data/Data.dat") //путь к файлу name.name
	if err != nil {
		return fmt.Errorf("ошибка чтения данных из файла: %w", err)
	}

	//генерация нового ID
	// Данные измерения, получить от пользователя
	measurementType := "temperature"
	measurementDate := "2025-03-28"
	id, nextNumber, err := manager.GenerateMeasurementID(measurementType, measurementDate)
	if err != nil {
		return fmt.Errorf("ошибка при генерации ID: %w", err)
	}

	// Вставляем каждую точку (x, y)
	for _, point := range result {
		_, err := manager.Db.Exec(models.InsertQuery, id, measurementType, measurementDate, nextNumber, title, point.X, point.Y)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// Получение данных из БД
func (manager *DBManager) GetDataFromDB(title models.Title) ([]models.Data, error) {
	// добавить ключи - конфиг?
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
	_, err := manager.Db.Exec(models.DeleteQuery)
	if err != nil {
		return fmt.Errorf("ошибка при удалении данных: %w", err)
	}
	fmt.Println("Данные успешно удалены")
	return nil
}
func (manager *DBManager) GenerateMeasurementID(measurementType, measurementDate string) (string, int, error) {
	var nextNumber int
	err := manager.Db.QueryRow(models.IDquery, measurementType, measurementDate).Scan(&nextNumber)
	if err != nil {
		return "", -1, fmt.Errorf("ошибка получения следующего номера измерения: %w", err)
	}

	id := fmt.Sprintf("%s_%s_%d", measurementType, measurementDate, nextNumber)
	return id, nextNumber, nil
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
