package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Pave1kor/Optics-characteristics/internal/app/repo"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		log.Println("Ошибка загрузки шаблона:", err)
		return
	}
	var experiment repo.DBInterface = repo.NewDBManager()

	err = experiment.ConnectToDb()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		log.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer experiment.Close()

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

func handlePostAction(w http.ResponseWriter, r *http.Request, experiment repo.DBInterface, action string) {
	switch action {
	case "load":
		err := experiment.AddDataToDB()
		if err != nil {
			http.Error(w, "Ошибка загрузки данных", http.StatusInternalServerError)
			log.Println("Ошибка загрузки данных:", err)
		}
		// Логика для загрузки данных
	case "delete":
		err := experiment.DeleteDataFromDB()
		if err != nil {
			http.Error(w, "Ошибка удаления данных", http.StatusInternalServerError)
			log.Println("Ошибка удаления данных:", err)
		}
		// Логика для удаления данных
	case "add":
		err := experiment.AddDataToDB()
		if err != nil {
			http.Error(w, "Ошибка добавления данных", http.StatusInternalServerError)
			log.Println("Ошибка добавления данных:", err)
		}
		// Логика для добавления данных
	case "drop":
		err := experiment.DropTable()
		if err != nil {
			http.Error(w, "Ошибка удаления таблицы", http.StatusInternalServerError)
			log.Println("Ошибка удаления таблицы:", err)
		}
		// Логика для удаления таблицы
	default:
		http.Error(w, "Некорректное действие", http.StatusBadRequest)
		log.Println("Некорректное действие:", action)
	}
}
