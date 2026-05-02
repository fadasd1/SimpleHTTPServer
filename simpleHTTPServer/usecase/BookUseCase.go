package usecase

import (
	"simpleHTTPServer/dto"
	"simpleHTTPServer/entity"
)

type BookRepository interface {
	GetByID(id string) (*entity.Book, error)
	CreateBook(request entity.Book) (string, error)
	GetBooks(filter dto.BookFilter) ([]entity.Book, error)
}

type BookUseCase struct {
	repo BookRepository
}

// validation should happen in this layer but since task is small in scope I skip
func NewBookUseCase(repo BookRepository) *BookUseCase {
	return &BookUseCase{repo: repo}
}

func (uc *BookUseCase) GetBookByID(id string) (*entity.Book, error) {
	return uc.repo.GetByID(id)
}

func (uc *BookUseCase) CreateBook(request dto.CreateBookRequest) (string, error) {
	return uc.repo.CreateBook(entity.Book{
		Title:    request.Title,
		AuthorID: request.AuthorID,
		Price:    request.Price,
		Stock:    request.Stock,
	})
}

func (uc *BookUseCase) GetBooks(filter dto.BookFilter) ([]entity.Book, error) {
	return uc.repo.GetBooks(filter)
}
