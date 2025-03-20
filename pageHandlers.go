package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	// "github.com/gorilla/sessions"
)

// var _:=seesions.NewCookieStore([]byte("secret-key"))
// handleHome обрабатывает главную страницу
func handleHome(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		log.Println("Ошибка загрузки шаблона:", err)
		return
	}
	db, err := connectToDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		log.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer db.Close()

	switch r.Method {
	case "GET":
		data, err := getDataFromDB(db)
		if err != nil {
			http.Error(w, "Ошибка получения данных из базы", http.StatusInternalServerError)
			log.Println("Ошибка получения данных из БД:", err)
			return
		}
		temp.Execute(w, data)

	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			log.Println("Ошибка обработки формы:", err)
			return
		}
		var err error
		action := r.FormValue("action")
		switch action {
		case "load":
			var data interface{}
			data, err = getDataFromDB(db)
			if err != nil {
				http.Error(w, "Ошибка получения данных из базы", http.StatusInternalServerError)
				log.Println("Ошибка получения данных из БД:", err)
				return
			}
			temp.Execute(w, data)
		case "delete":
			err = deleteDataFromDB(db)
			if err != nil {
				http.Error(w, "Ошибка удаления данных", http.StatusInternalServerError)
				log.Println("Ошибка удаления данных:", err)
				return
			}
			temp.Execute(w, nil)
		case "add":
			err = addDataToDB(db)
			if err != nil {
				http.Error(w, "Ошибка добавления данных", http.StatusInternalServerError)
				log.Println("Ошибка добавления данных:", err)
				return
			}
			temp.Execute(w, nil)
		default:
			err = fmt.Errorf("некорректное действие")
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Ошибка:", err)
			return
		}
	}
}

// handleAbout handles the about page
func handleAbout(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	temp.Execute(w, nil)
}

// handleContact handles the contact page
func handleContact(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/contact.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	temp.Execute(w, nil)
}
