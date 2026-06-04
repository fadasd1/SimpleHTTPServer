// usecase/book_usecase_test.go
package usecase

import (
	"context"
	"errors"
	"fmt"
	"simpleHTTPServer/dto"
	"testing"

	"simpleHTTPServer/apperror"
	"simpleHTTPServer/entity"
	"simpleHTTPServer/usecase/mocks"

	"go.uber.org/mock/gomock"
)

func TestGetBookByID_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	expected := &entity.Book{ID: 1, Title: "Test Book", AuthorID: 1, Price: 10.99, Stock: 5}

	// tell the mock: when GetByID is called with any context and id "1", return expected
	mockRepo.EXPECT().
		GetByID(gomock.Any(), "1").
		Return(expected, nil)

	uc := NewBookUseCase(mockRepo)
	book, err := uc.GetBookByID(context.Background(), "1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if *book != *expected {
		t.Errorf("expected %v, got %v", expected, book)
	}
}

func TestGetBookByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)

	mockRepo.EXPECT().
		GetByID(gomock.Any(), "99").
		Return(nil, fmt.Errorf("book with id 99: %w", apperror.ErrNotFound))

	uc := NewBookUseCase(mockRepo)
	_, err := uc.GetBookByID(context.Background(), "99")

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCreateBook_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	request := dto.CreateBookRequest{Title: "New Book", AuthorID: 1, Price: 19.99, Stock: 10}

	mockRepo.EXPECT().
		CreateBook(gomock.Any(), request).
		Return("book created successfully", nil)

	uc := NewBookUseCase(mockRepo)
	msg, err := uc.CreateBook(context.Background(), request)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if msg != "book created successfully" {
		t.Errorf("expected success message, got %v", msg)
	}
}

func TestCreateBook_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	request := dto.CreateBookRequest{Title: "New Book", AuthorID: 1, Price: 19.99, Stock: 10}

	mockRepo.EXPECT().
		CreateBook(gomock.Any(), gomock.Any()).
		Return("", fmt.Errorf("CreateBook query: %w", apperror.ErrInternal))

	uc := NewBookUseCase(mockRepo)
	_, err := uc.CreateBook(context.Background(), request)

	if !errors.Is(err, apperror.ErrInternal) {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func TestGetBooks_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	expected := []entity.Book{
		{ID: 1, Title: "Book One", AuthorID: 1, Price: 9.99, Stock: 3},
		{ID: 2, Title: "Book Two", AuthorID: 2, Price: 14.99, Stock: 7},
	}

	mockRepo.EXPECT().
		GetBooks(gomock.Any(), dto.BookFilter{}).
		Return(expected, nil)

	uc := NewBookUseCase(mockRepo)
	books, err := uc.GetBooks(context.Background(), dto.BookFilter{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(books) != 2 {
		t.Errorf("expected 2 books, got %d", len(books))
	}
}

func TestGetBooks_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)

	mockRepo.EXPECT().
		GetBooks(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("GetBooks query: %w", apperror.ErrInternal))

	uc := NewBookUseCase(mockRepo)
	_, err := uc.GetBooks(context.Background(), dto.BookFilter{})

	if !errors.Is(err, apperror.ErrInternal) {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}
