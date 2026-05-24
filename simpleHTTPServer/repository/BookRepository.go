// repository/book_repository.go
package repository

import (
	"context"
	"errors"
	"fmt"
	"simpleHTTPServer/apperror"
	"simpleHTTPServer/dto"
	"simpleHTTPServer/entity"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

type BookRepository struct {
	db DB
}

func NewBookRepository(db DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetByID(ctx context.Context, id string) (*entity.Book, error) {
	book := &entity.Book{}
	err := r.db.QueryRow(ctx, `SELECT id, title, author_id, price, stock FROM books WHERE id = $1`, id).
		Scan(&book.ID, &book.Title, &book.AuthorID, &book.Price, &book.Stock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("book with id %s: %w", id, apperror.ErrNotFound)
		}
		return nil, fmt.Errorf("GetByID query: %w", apperror.ErrInternal)
	}
	return book, nil
}

func (r *BookRepository) CreateBook(ctx context.Context, request dto.CreateBookRequest) (string, error) {
	_, err := r.db.Exec(ctx, `
        INSERT INTO books (title, author_id, price, stock)
        VALUES ($1, $2, $3, $4)
    `, request.Title, request.AuthorID, request.Price, request.Stock)
	if err != nil {
		return "", fmt.Errorf("CreateBook query: %w", apperror.ErrInternal)
	}
	return "book created successfully", nil
}

func (r *BookRepository) GetBooks(ctx context.Context, filter dto.BookFilter) ([]entity.Book, error) {
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

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetBooks query: %w", apperror.ErrInternal)
	}
	defer rows.Close()

	books := make([]entity.Book, 0)
	for rows.Next() {
		var b entity.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.AuthorID, &b.Price, &b.Stock); err != nil {
			return nil, fmt.Errorf("GetBooks scan: %w", apperror.ErrInternal)
		}
		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetBooks rows: %w", apperror.ErrInternal)
	}

	return books, nil
}
