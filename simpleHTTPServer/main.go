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
	http.HandleFunc("/books/get", getBooksHandler)

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

func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	authorID := r.URL.Query().Get("author_id")
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")
	inStock := r.URL.Query().Get("in_stock")

	page := 1
	limit := 10

	if pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &page)
	}
	if limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	offset := (page - 1) * limit

	query := `
		SELECT id, title, author_id, price, stock
		FROM books
		WHERE 1=1
	`

	args := make([]any, 0)
	argID := 1

	if authorID != "" {
		query += fmt.Sprintf(" AND author_id = $%d", argID)
		args = append(args, authorID)
		argID++
	}

	if minPrice != "" {
		query += fmt.Sprintf(" AND price >= $%d", argID)
		args = append(args, minPrice)
		argID++
	}

	if maxPrice != "" {
		query += fmt.Sprintf(" AND price <= $%d", argID)
		args = append(args, maxPrice)
		argID++
	}

	if inStock == "true" {
		query += " AND stock > 0"
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, authorID, stock int
		var title string
		var price float64

		err := rows.Scan(&id, &title, &authorID, &price, &stock)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "ID:%d | %s | Author:%d | %.2f | Stock:%d\n",
			id, title, authorID, price, stock)
	}
}
