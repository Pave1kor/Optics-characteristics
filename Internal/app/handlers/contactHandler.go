package handlers

import (
	"html/template"
	"net/http"
)

func HandlerContact(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/user/contact.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	temp.Execute(w, nil)
}
