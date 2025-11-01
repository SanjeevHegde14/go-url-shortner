package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3" // Blank import for the SQLite3 driver
)

var db *sql.DB

type ShortenRequest struct {
	URL string `json:"url"` // JSON struct tag for unmarshaling
}

func generateShortCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./urls.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS urls (
		"short_code" TEXT PRIMARY KEY,
		"long_url" TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %q", err)
	}

	log.Println("Database initialized successfully")
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortCode := generateShortCode(6)

	insertSQL := `INSERT INTO urls(short_code, long_url) VALUES (?, ?)`
	_, err = db.Exec(insertSQL, shortCode, req.URL)
	if err != nil {
		log.Printf("Error inserting into database: %v\n", err)
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"short_url": fmt.Sprintf("http://localhost:8080/%s", shortCode),
		"long_url":  req.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	shortCode := r.URL.Path[1:]

	if shortCode == "" {
		fmt.Fprintln(w, "Go URL Shortener is live. POST to /shorten to create.")
		return
	}

	var longURL string
	querySQL := `SELECT long_url FROM urls WHERE short_code = ?`
	row := db.QueryRow(querySQL, shortCode)

	err := row.Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Short URL not found", http.StatusNotFound)
		} else {
			log.Printf("Error querying database: %v\n", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", handleRedirect)
	http.HandleFunc("/shorten", handleShorten)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}