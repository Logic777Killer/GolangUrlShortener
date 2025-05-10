package shortener

import (
	"crypto/md5"
	"fmt"
	"net/http"
)

type URLShortener struct {
	store map[string]string
}

func New() *URLShortener {
	return &URLShortener{
		store: make(map[string]string),
	}
}

func (s *URLShortener) GenerateShortURL(original string) string {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(original)))[:5]
	s.store[hash] = original
	return hash
}

func (s *URLShortener) RedirectHandler(w http.ResponseWriter, r *http.Request, shortURL string) {
	original, exists := s.store[shortURL]
	if !exists {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, original, http.StatusMovedPermanently)
}
