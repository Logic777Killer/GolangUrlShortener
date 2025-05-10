package main

import (
	"fmt"
	"log"
	"net/http"

	"GolangUrlShortenerWeb/shortener"
	"GolangUrlShortenerWeb/web"

	"github.com/gorilla/mux"
)

func main() {
	// Инициализация компонентов
	shortener := shortener.New()
	handlers := web.NewHandlers(shortener)

	// Настройка маршрутов
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Home).Methods("GET")
	r.HandleFunc("/", handlers.HandleShorten).Methods("POST")

	// Обработчик для редиректов
	r.HandleFunc("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortener.RedirectHandler(w, r, vars["shortURL"])
	}).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
