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

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next(w, r)
	}
}
