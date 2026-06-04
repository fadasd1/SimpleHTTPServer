package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"simpleHTTPServer/apperror"
	"simpleHTTPServer/dto"
	"simpleHTTPServer/entity"
)

type BookUseCase interface {
	GetBookByID(ctx context.Context, id string) (*entity.Book, error)
	CreateBook(ctx context.Context, request dto.CreateBookRequest) (string, error)
	GetBooks(ctx context.Context, filter dto.BookFilter) ([]entity.Book, error)
}

type BookHandler struct {
	useCase BookUseCase
}

func NewBookHandler(uc BookUseCase) *BookHandler {
	return &BookHandler{useCase: uc}
}

func httpError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperror.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, apperror.ErrAlreadyExists):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, apperror.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *BookHandler) GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	var filter dto.BookFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	books, err := h.useCase.GetBooks(r.Context(), filter)
	if err != nil {
		httpError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	for _, b := range books {
		fmt.Fprintf(w, "ID:%d | %s | Author:%d | Price:%.2f | Stock:%d\n",
			b.ID, b.Title, b.AuthorID, b.Price, b.Stock)
	}
}

func (h *BookHandler) GetBookByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	book, err := h.useCase.GetBookByID(r.Context(), id)
	if err != nil {
		httpError(w, err)
		return
	}

	fmt.Fprintf(w, "ID: %d | %s | Price: %.2f | Stock: %d\n",
		book.ID, book.Title, book.Price, book.Stock)
}

func (h *BookHandler) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	msg, err := h.useCase.CreateBook(r.Context(), request)
	if err != nil {
		httpError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, msg)
}
