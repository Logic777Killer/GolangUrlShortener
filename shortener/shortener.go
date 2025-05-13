package shortener

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type URLShortener struct {
	db *sql.DB
}

// New инициализирует подключение к PostgreSQL
func New(dsn string) (*URLShortener, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверка подключения
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %v", err)
	}

	return &URLShortener{db: db}, nil
}

// GenerateShortURL создает или возвращает существующий короткий код
func (s *URLShortener) GenerateShortURL(originalURL string) (string, error) {
	// Проверка существующей записи
	var existingCode string
	err := s.db.QueryRow(
		"SELECT short_code FROM urls WHERE original_url = $1",
		originalURL,
	).Scan(&existingCode)

	if err == nil {
		return existingCode, nil
	}

	// Генерация нового кода
	hash := md5.Sum([]byte(originalURL))
	shortCode := fmt.Sprintf("%x", hash)[:5]

	// Вставка новой записи
	_, err = s.db.Exec(
		"INSERT INTO urls (original_url, short_code) VALUES ($1, $2)",
		originalURL,
		shortCode,
	)

	if err != nil {
		return "", fmt.Errorf("insert failed: %v", err)
	}

	return shortCode, nil
}

// RedirectHandler обрабатывает редирект по короткому коду
func (s *URLShortener) RedirectHandler(w http.ResponseWriter, r *http.Request, shortCode string) {
	var originalURL string
	err := s.db.QueryRow(
		"SELECT original_url FROM urls WHERE short_code = $1",
		shortCode,
	).Scan(&originalURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
