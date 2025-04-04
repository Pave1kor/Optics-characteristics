package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Кэшируем шаблоны при инициализации пакета
var aboutTemplates = template.Must(template.ParseFiles(
	filepath.Join("templates", "base.html"),
	filepath.Join("templates", "user", "about.html"),
))

func HandlerAbout(w http.ResponseWriter, r *http.Request) {
	// Выполняем шаблон с обработкой ошибок
	err := aboutTemplates.ExecuteTemplate(w, "base.html", nil)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
