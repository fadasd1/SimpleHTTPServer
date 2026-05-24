package repository

import (
	"context"
	"errors"
	"simpleHTTPServer/apperror"
	"simpleHTTPServer/dto"
	"simpleHTTPServer/entity"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
)

func TestGetByID_HappyPath(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	expected := &entity.Book{ID: 1, Title: "Test Book", AuthorID: 1, Price: 10.99, Stock: 5}

	mock.ExpectQuery(`SELECT id, title, author_id, price, stock FROM books WHERE id = \$1`).
		WithArgs("1").
		WillReturnRows(pgxmock.NewRows([]string{"id", "title", "author_id", "price", "stock"}).
			AddRow(expected.ID, expected.Title, expected.AuthorID, expected.Price, expected.Stock))

	repo := NewBookRepository(mock)

	book, err := repo.GetByID(context.Background(), "1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if *book != *expected {
		t.Errorf("expected %v, got %v", expected, book)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetByID_BookNotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectQuery(`SELECT id, title, author_id, price, stock FROM books WHERE id = \$1`).
		WithArgs("1").
		WillReturnError(pgx.ErrNoRows)

	repo := NewBookRepository(mock)

	book, err := repo.GetByID(context.Background(), "1")
	if book != nil {
		t.Fatalf("expected no book, got %v", book)
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected %v, got %v", apperror.ErrNotFound, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestCreateBook_HappyPath(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	bookRequest := dto.CreateBookRequest{Title: "book", AuthorID: 1, Price: 1, Stock: 1}

	mock.ExpectExec(`INSERT INTO books`).
		WithArgs(bookRequest.Title, bookRequest.AuthorID, bookRequest.Price, bookRequest.Stock).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	repo := NewBookRepository(mock)
	msg, err := repo.CreateBook(context.Background(), bookRequest)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if msg != "book created successfully" {
		t.Errorf("success message not as designated")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestCreateBook_InternalError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	bookRequest := dto.CreateBookRequest{Title: "book", AuthorID: 1, Price: 1, Stock: 1}

	mock.ExpectExec(`INSERT INTO books`).
		WithArgs(bookRequest.Title, bookRequest.AuthorID, bookRequest.Price, bookRequest.Stock).
		WillReturnError(errors.New("db connection lost"))

	repo := NewBookRepository(mock)
	msg, err := repo.CreateBook(context.Background(), bookRequest)

	if msg != "" {
		t.Fatalf("expected empty msg, got %s", msg)
	}

	if !errors.Is(err, apperror.ErrInternal) {
		t.Errorf("expected %v, got %v", apperror.ErrInternal, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetBooks_HappyPath(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"id", "title", "author_id", "price", "stock"}).
		AddRow(1, "Book One", 1, 9.99, 3).
		AddRow(2, "Book Two", 2, 14.99, 7)
	mock.ExpectQuery(`SELECT id, title, author_id, price, stock FROM books WHERE 1=1 LIMIT \$1 OFFSET \$2`).
		WithArgs(10, 0).
		WillReturnRows(rows)

	repo := NewBookRepository(mock)
	books, err := repo.GetBooks(context.Background(), dto.BookFilter{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(books) != 2 {
		t.Errorf("expected 2 books, received %d", len(books))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetBooks_InternalError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectQuery(`SELECT id, title, author_id, price, stock FROM books WHERE 1=1 LIMIT \$1 OFFSET \$2`).
		WithArgs(10, 0).
		WillReturnError(errors.New("db connection lost"))

	repo := NewBookRepository(mock)
	books, err := repo.GetBooks(context.Background(), dto.BookFilter{})

	if books != nil {
		t.Fatalf("expected no books, got %v", books)
	}

	if !errors.Is(err, apperror.ErrInternal) {
		t.Errorf("expected %v, got %v", apperror.ErrInternal, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
