package main

import (
	"html/template"
	"net/http"
	// "github.com/gorilla/sessions"
)

// var _:=seesions.NewCookieStore([]byte("secret-key"))
// handleHome handles the home page
func handleHome(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		r.ParseForm()
		db, err := ConnectToDB()
		if err != nil {
			http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		temp.Execute(w, nil)
	}

	temp.Execute(w, nil)
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
