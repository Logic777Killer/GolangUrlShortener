package web

import (
	"html/template"
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

func NewHandlers(shortener *shortener.URLShortener) *Handlers {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("web", "templates", "base.html"),
		filepath.Join("web", "templates", "styles.css"),
	))

	return &Handlers{
		shortener: shortener,
		tmpl:      tmpl,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, PageData{})
}

func (h *Handlers) HandleShorten(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.Form.Get("url")
	data := PageData{OriginalURL: url}

	if url == "" {
		data.Error = "URL не может быть пустым"
		h.renderTemplate(w, data)
		return
	}

	shortURL := h.shortener.GenerateShortURL(url)
	data.ShortURL = shortURL
	h.renderTemplate(w, data)
}

func (h *Handlers) renderTemplate(w http.ResponseWriter, data PageData) {
	if err := h.tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
	}
}
