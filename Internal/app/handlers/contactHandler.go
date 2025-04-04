package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Кэшируем шаблоны при инициализации
var contactTemplates = template.Must(template.ParseFiles(
	filepath.Join("templates", "base.html"),
	filepath.Join("templates", "user", "contact.html"),
))

func HandlerContact(w http.ResponseWriter, r *http.Request) {
	// Рендеринг с обработкой ошибок
	err := contactTemplates.ExecuteTemplate(w, "base.html", nil)
	if err != nil {
		log.Printf("Ошибка рендеринга шаблона: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
