package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Pave1kor/Optics-characteristics/internal/app/models"
	"github.com/Pave1kor/Optics-characteristics/internal/app/repo"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		log.Println("Ошибка загрузки шаблона:", err)
		return
	}
	dbManager := &models.DBManager{}
	experiment := repo.NewDBWrapper(dbManager) // Создаём обёртку
	err = experiment.ConnectToDb()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		log.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer experiment.DBManager.Db.Close()

	switch r.Method {
	case "GET":
		files, err := experiment.GetListOfFiles()
		if err != nil {
			http.Error(w, "Ошибка получения списка файлов", http.StatusInternalServerError)
			log.Println("Ошибка получения списка файлов:", err)
			return
		}
		if len(files) == 0 {
			log.Println("Список файлов пуст")
			temp.Execute(w, []string{})
			return
		}
		if err := temp.Execute(w, files); err != nil {
			http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
			log.Println("Ошибка рендеринга шаблона:", err)
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			log.Println("Ошибка обработки формы:", err)
			return
		}
		action := r.FormValue("action")
		handlePostAction(w, r, experiment, action)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		log.Println("Метод не поддерживается:", r.Method)
	}
}

func handlePostAction(w http.ResponseWriter, r *http.Request, experiment *repo.DBWrapper, action string) {
	switch action {
	case "load":
		// Логика для загрузки данных
	case "delete":
		// Логика для удаления данных
	case "add":
		// Логика для добавления данных
	case "drop":
		// Логика для удаления таблицы
	case "create":
		// Логика для создания таблицы
	default:
		http.Error(w, "Некорректное действие", http.StatusBadRequest)
		log.Println("Некорректное действие:", action)
	}
}
