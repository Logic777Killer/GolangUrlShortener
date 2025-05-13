package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"GolangUrlShortenerWeb/shortener"
	"GolangUrlShortenerWeb/web"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {
	// Конфигурация PostgreSQL
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Формирование DSN из переменных окружения
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Инициализация shortener
	shortenerInstance, err := shortener.New(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize shortener: %v", err)
	}

	// Создание обработчиков
	handlers := web.NewHandlers(shortenerInstance)

	// Настройка маршрутизатора
	r := mux.NewRouter()

	// Статические файлы
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	// Маршруты
	r.HandleFunc("/", handlers.Home).Methods("GET")
	r.HandleFunc("/", handlers.HandleShorten).Methods("POST")
	r.HandleFunc("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortenerInstance.RedirectHandler(w, r, vars["shortURL"])
	}).Methods("GET")

	// Запуск сервера
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
