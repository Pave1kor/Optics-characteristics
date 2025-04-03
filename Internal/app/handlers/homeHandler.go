package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Pave1kor/Optics-characteristics/internal/app/models"
	"github.com/Pave1kor/Optics-characteristics/internal/app/repo"
)

// Кэшируем шаблон при инициализации
var homeTemplate = template.Must(template.ParseFiles("templates/index.html"))

func HandleHome(w http.ResponseWriter, r *http.Request) {
	experiment := repo.NewDBManager()
	if err := experiment.ConnectToDb(); err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		log.Println("DB connection error:", err)
		return
	}
	defer experiment.Close()

	switch r.Method {
	case http.MethodGet:
		handleGetRequest(w, experiment)
	case http.MethodPost:
		handlePostRequest(w, r, experiment)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetRequest(w http.ResponseWriter, db *repo.DBManager) {
	// Проверяем и инициализируем БД только при первом запросе
	if err := initializeDatabase(db); err != nil {
		http.Error(w, "Database initialization error", http.StatusInternalServerError)
		log.Println("DB init error:", err)
		return
	}
	result, err := db.GetDataFromDB()
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		log.Println("DB query error:", err)
		return
	}
	// Передаем данные в шаблон
	templateData := struct {
		Title string
		Data  []models.Data // Замените на вашу модель данных
	}{
		Title: "Оптические характеристики",
		Data:  result,
	}

	if err := homeTemplate.Execute(w, templateData); err != nil {
		log.Println("Template execution error:", err)
	}
}

func initializeDatabase(db *repo.DBManager) error {
	// Создаем таблицы, если они не существуют
	if err := db.CheckTables(models.TableName); err != nil {
		return fmt.Errorf("table check failed: %w", err)
	}

	exists, err := db.CheckDataInTables()
	if err != nil {
		return fmt.Errorf("data check failed: %w", err)
	}
	//Добавляем данные в таблицу, если они отсутствуют
	if !exists {
		if err := db.AddDataToDB(); err != nil {
			return fmt.Errorf("data insertion failed: %w", err)
		}
	}
	return nil
}

func handlePostRequest(w http.ResponseWriter, r *http.Request, experiment *repo.DBManager) {
	action := r.FormValue("action")
	switch action {
	case "load":
		err := experiment.AddDataToDB()
		if err != nil {
			http.Error(w, "Ошибка загрузки данных", http.StatusInternalServerError)
			log.Println("Ошибка загрузки данных:", err)
		}
		// Логика для загрузки данных
	case "add":
		err := experiment.AddDataToDB()
		if err != nil {
			http.Error(w, "Ошибка добавления данных", http.StatusInternalServerError)
			log.Println("Ошибка добавления данных:", err)
		}
		// Логика для добавления данных
	default:
		http.Error(w, "Некорректное действие", http.StatusBadRequest)
		log.Println("Некорректное действие:", action)
	}
}
