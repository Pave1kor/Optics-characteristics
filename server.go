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

// Add data to baseData
func (manager *DBManager) addDataToDB(name string) (Title, error) {
	//Load data - сделать универсальным
	result, title, err := readDataFromFile("data/Data.dat") //путь к файлу name.name
	if err != nil {
		return Title{}, fmt.Errorf("ошибка чтения данных из файла: %w", err)
	}
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" ("%s" FLOAT PRIMARY KEY, "%s" FLOAT NOT NULL)`,
		name, title.X, title.Y)

	_, err = manager.db.Exec(query)
	if err != nil {
		return Title{}, fmt.Errorf("ошибка при создании таблицы: %w", err)
	}
	fmt.Println("Таблица успешно создана")

	query = fmt.Sprintf(`INSERT INTO %s ("%s", "%s") VALUES (&1, &2)ON CONFLICT (id) DO UPDATE SET "%s" = EXCLUDED."%s", "%s" = EXCLUDED."%s";`,
		name, title.X, title.Y, title.X, title.X, title.Y, title.Y)

	// Вставка данных
	for _, data := range result {
		_, err := manager.db.Exec(query, data.X, data.Y)
		if err != nil {
			return Title{}, fmt.Errorf("ошибка при вставке данных: %w", err)
		}
	}
	fmt.Println("Данные успешно добавлены")
	return title, nil
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
