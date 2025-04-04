package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Pave1kor/Optics-characteristics/internal/app/handlers"
	"github.com/spf13/viper"
)

func main() {
	// Настройка конфига
	viper.SetDefault("PORT", "8080")
	port := viper.GetString("PORT")

	// Инициализация сервера
	server := &http.Server{
		Addr:    ":" + port,
		Handler: setupRoutes(),
	}

	// Graceful shutdown
	go func() {
		log.Printf("Сервер запущен на http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при завершении сервера: %v", err)
	}
	log.Println("Сервер корректно остановлен")
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	// Обработка статических файлов
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", loggingMiddleware(corsMiddleware(handlers.HandleHome)))
	mux.HandleFunc("/about", loggingMiddleware(handlers.HandlerAbout))
	mux.HandleFunc("/contact", loggingMiddleware(handlers.HandlerContact))
	return mux
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Запрос: %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Ответ: %s %s за %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// This function creates a middleware that sets the CORS headers for an HTTP request
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// This function returns a new http.HandlerFunc that sets the CORS headers and calls the next handler
	return func(w http.ResponseWriter, r *http.Request) {
		// Set the Access-Control-Allow-Origin header to allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Set the Access-Control-Allow-Methods header to allow GET, POST, and OPTIONS requests
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		// Set the Access-Control-Allow-Headers header to allow the Content-Type header
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Call the next handler
		next(w, r)
	}
}
