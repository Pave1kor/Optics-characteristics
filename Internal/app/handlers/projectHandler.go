package main

import (
	"html/template"
	"net/http"
)

func handlerProject(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/project.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	temp.Execute(w, nil)
}
