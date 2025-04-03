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
	AddDataToDB() error
	CheckTables(tableName string) error
	CheckDataInTables(tableName string) (bool, error)
	ConnectToDb() error
	Close() error
	DeleteDataFromDB() error
	DropTable() error
	GenerateMeasurementID(measurementType, measurementDate string) (string, int, error)
	GetDataFromDB() ([]models.Data, error)
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

// Check tables
func (manager *DBManager) CheckTables(tableName string) error {
	exists, err := TablesExist(manager.Db, tableName)
	if err != nil {
		return fmt.Errorf("ошибка проверки таблицы: %w", err)
	}
	if exists {
		return nil
	}
	_, err = manager.Db.Exec(fmt.Sprintf(models.CreateTableQuery, tableName))
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %v", err)
	}

	log.Printf("Таблица %s создана", tableName)
	return nil
}

func TablesExist(dB *sql.DB, tableName string) (bool, error) {
	// Проверяем существование таблицы
	var exists bool
	err := dB.QueryRow(models.CheckTableQuery, tableName).Scan(&exists)
	if err != nil || !exists {
		return false, err
	}
	return true, nil
}

// Check data in table
func (manager *DBManager) CheckDataInTables() (bool, error) {
	var exists bool
	if err := manager.Db.QueryRow(models.CheckDataInTablesQuery).Scan(&exists); err != nil {
		return false, fmt.Errorf("ошибка проверки данных в таблице: %w", err)
	}
	return exists, nil
}

// Get list of files
// func (manager *DBManager) GetListOfTables() ([]models.TableName, error) {
// 	rows, err := manager.Db.Query(models.GetListOfTablesQuery)
// 	list := []models.TableName{}
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	// Вывод списка таблиц
// 	fmt.Println("Список таблиц:")
// 	for rows.Next() {
// 		var tableName models.TableName
// 		if err := rows.Scan(&tableName); err != nil {
// 			log.Fatal(err)
// 		}
// 		list = append(list, tableName)
// 	}
// 	return list, nil
// }

// Add data to baseData
func (manager *DBManager) AddDataToDB() error {
	//Load data - сделать универсальным
	result, _, err := read.ReadDataFromFile("data/Data.dat") //путь к файлу name.name
	if err != nil {
		return fmt.Errorf("ошибка чтения данных из файла: %w", err)
	}
	// Переименование шапки таблицы - добавить переименование столбцов
	// _, err = manager.Db.Exec(fmt.Sprint(models.RenameTableQuery, title.X, title.Y))
	// if err != nil {
	// 	return fmt.Errorf("ошибка при переименовании столбцов: %w", err)
	// }

	// //генерация ID - таблицы (уникальное имя?)
	// // Данные измерения, получить от пользователя
	// measurementType := "temperature"
	// measurementDate := "2025-03-28"
	// id, err := manager.GenerateMeasurementID(measurementType, measurementDate)
	// if err != nil {
	// 	return fmt.Errorf("ошибка при генерации ID: %w", err)
	// }

	// Вставляем каждую точку (x, y)

	for _, point := range result {
		_, err := manager.Db.Exec(models.InsertQuery, point.X, point.Y)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// Получение данных из БД
func (manager *DBManager) GetDataFromDB() ([]models.Data, error) {
	// добавить ключи - конфиг?
	rows, err := manager.Db.Query(models.SelectDataQuery)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе данных: %w", err)
	}
	defer rows.Close()

	var results []models.Data
	for rows.Next() {
		var data models.Data
		if err := rows.Scan(&data.Id, &data.X, &data.Y); err != nil {
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
//
//	func (manager *DBManager) DeleteDataFromDB() error {
//		//добавить ключи
//		// query := fmt.Sprintf(`DELETE FROM %s`,
//		// 	name)
//		_, err := manager.Db.Exec(models.DeleteQuery)
//		if err != nil {
//			return fmt.Errorf("ошибка при удалении данных: %w", err)
//		}
//		fmt.Println("Данные успешно удалены")
//		return nil
//	}
// func (manager *DBManager) GenerateMeasurementID(measurementType, measurementDate string) (string, error) {
// 	var nextNumber int
// 	err := manager.Db.QueryRow(models.IDquery, measurementType, measurementDate).Scan(&nextNumber)
// 	if err != nil {
// 		return "", fmt.Errorf("ошибка получения следующего номера измерения: %w", err)
// 	}

// 	id := fmt.Sprintf("%s_%s_%d", measurementType, measurementDate, nextNumber)
// 	return id, nil
// }

// Удаление таблицы
// func (manager *DBManager) DropTable() error { //удалить по запросу
// 	query := fmt.Sprintf(`DROP TABLE IF EXISTS measurement`)
// 	_, err := manager.Db.Exec(query)
// 	if err != nil {
// 		return fmt.Errorf("ошибка при удалении таблицы: %w", err)
// 	}
// 	fmt.Println("Таблица успешно удалена")
// 	return nil
// }

func (manager *DBManager) Close() error {
	return manager.Db.Close()
}
