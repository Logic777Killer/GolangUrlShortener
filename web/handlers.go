package web

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"GolangUrlShortenerWeb/shortener"
)

type PageData struct {
	OriginalURL string
	ShortURL    string
	Error       string
}

type Handlers struct {
	shortener *shortener.URLShortener
	tmpl      *template.Template
}

func NewHandlers(s *shortener.URLShortener) *Handlers {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("web", "templates", "base.html"),
		filepath.Join("web", "templates", "styles.css"),
	))
	return &Handlers{
		shortener: s,
		tmpl:      tmpl,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, PageData{})
}

func (h *Handlers) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderTemplate(w, PageData{Error: "Invalid form data"})
		return
	}

	originalURL := r.Form.Get("url")
	if originalURL == "" {
		h.renderTemplate(w, PageData{Error: "URL cannot be empty"})
		return
	}

	shortCode, err := h.shortener.GenerateShortURL(originalURL)
	if err != nil {
		log.Printf("Ошибка при создании короткой ссылки: %v", err) // Добавлено
		h.renderTemplate(w, PageData{
			OriginalURL: originalURL,
			Error:       "Failed to create short URL",
		})
		return
	}
	
	h.renderTemplate(w, PageData{
		OriginalURL: originalURL,
		ShortURL:    shortCode,
	})
}

func (h *Handlers) renderTemplate(w http.ResponseWriter, data PageData) {
	err := h.tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}
