package usecase

import (
	"context"
	"fmt"
	"simpleHTTPServer/dto"
	"simpleHTTPServer/entity"
)

type BookRepository interface {
	GetByID(ctx context.Context, id string) (*entity.Book, error)
	CreateBook(ctx context.Context, request entity.Book) (string, error)
	GetBooks(ctx context.Context, filter dto.BookFilter) ([]entity.Book, error)
}

type BookUseCase struct {
	repo BookRepository
}

func NewBookUseCase(repo BookRepository) *BookUseCase {
	return &BookUseCase{repo: repo}
}

func (uc *BookUseCase) GetBookByID(ctx context.Context, id string) (*entity.Book, error) {
	book, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetBookByID: %w", err)
	}
	return book, nil
}

func (uc *BookUseCase) CreateBook(ctx context.Context, request dto.CreateBookRequest) (string, error) {
	msg, err := uc.repo.CreateBook(ctx, entity.Book{
		Title:    request.Title,
		AuthorID: request.AuthorID,
		Price:    request.Price,
		Stock:    request.Stock,
	})
	if err != nil {
		return "", fmt.Errorf("CreateBook: %w", err)
	}
	return msg, nil
}

func (uc *BookUseCase) GetBooks(ctx context.Context, filter dto.BookFilter) ([]entity.Book, error) {
	books, err := uc.repo.GetBooks(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("GetBooks: %w", err)
	}
	return books, nil
}
