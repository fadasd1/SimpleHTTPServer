package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func main() {

	var err error

	connStr := "postgres://fadasd:1337@dummyDB:5432/dummyDB?sslmode=disable"

	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/books/id/", getBookByIDHandler)
	http.HandleFunc("/books/create", createBookHandler)

	log.Println("API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getBookByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	idStr := strings.TrimPrefix(r.URL.Path, "/books/id/")

	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	var bookID int
	var title string
	var authorID int
	var price float64
	var stock int

	err := db.QueryRow(`
		SELECT id, title, author_id, price, stock
		FROM books
		WHERE id = $1
	`, idStr).Scan(&bookID, &title, &authorID, &price, &stock)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "ID: %d | %s | Price: %.2f | Stock: %d\n",
		bookID, title, price, stock)
}

func createBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	title := r.URL.Query().Get("title")
	authorID := r.URL.Query().Get("author_id")
	price := r.URL.Query().Get("price")
	stock := r.URL.Query().Get("stock")

	_, err := db.Exec(`
		INSERT INTO books (title, author_id, price, stock)
		VALUES ($1, $2, $3, $4)
	`, title, authorID, price, stock)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "book created")
}
