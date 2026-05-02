package repository

import (
	"database/sql"
	"fmt"
	"simpleHTTPServer/dto"
	"simpleHTTPServer/entity"
	"strconv"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetByID(id string) (*entity.Book, error) {
	var b entity.Book
	err := r.db.QueryRow(`SELECT id, title, author_id, price, stock FROM books WHERE id = $1`, id).
		Scan(&b.ID, &b.Title, &b.AuthorID, &b.Price, &b.Stock)
	if err != nil {
		return nil, err
	}
	return &b, err
}

func (r *BookRepository) CreateBook(request entity.Book) (string, error) {

	_, err := r.db.Exec(`
INSERT INTO books (title, author_id, price, stock)
VALUES ($1, $2, $3, $4)
`, request.Title, request.AuthorID, request.Price, request.Stock)

	if err != nil {
		return "error", err
	}
	return "success", err
}

func (r *BookRepository) GetBooks(filter dto.BookFilter) ([]entity.Book, error) {

	query := `SELECT id, title, author_id, price, stock FROM books WHERE 1=1`

	args := make([]any, 0)
	argID := 1

	if filter.AuthorID != "" {
		query += fmt.Sprintf(" AND author_id = $%d", argID)
		args = append(args, filter.AuthorID)
		argID++
	}

	if filter.MinPrice != "" {
		query += fmt.Sprintf(" AND price >= $%d", argID)
		args = append(args, filter.MinPrice)
		argID++
	}

	if filter.InStock == "true" {
		query += " AND stock > 0"
	}

	limit := 10
	offset := 0

	if filter.Limit != "" {
		fmt.Sscanf(filter.Limit, "%d", &limit)
	}

	if filter.Page != "" {
		page, _ := strconv.Atoi(filter.Page)
		offset = (page - 1) * limit
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...) // The "..." unpacks the slice into separate arguments
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	books := make([]entity.Book, 0)

	for rows.Next() {
		var id, authorID, stock int
		var title string
		var price float64

		err := rows.Scan(&id, &title, &authorID, &price, &stock)
		if err != nil {
			return nil, err
		}
		book := entity.Book{ID: id, Title: title, AuthorID: authorID, Price: price, Stock: stock}

		books = append(books, book)

	}

	return books, nil
}
