package controller

import (
	"fmt"
	"net/http"
	"simpleHTTPServer/dto"
	"simpleHTTPServer/usecase"
	"strconv"
	"strings"
)

type BookHandler struct {
	useCase *usecase.BookUseCase
}

func NewBookHandler(uc *usecase.BookUseCase) *BookHandler {
	return &BookHandler{useCase: uc}
}

func (h *BookHandler) GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	filter := dto.BookFilter{
		Page:     r.URL.Query().Get("page"),
		Limit:    r.URL.Query().Get("limit"),
		AuthorID: r.URL.Query().Get("author_id"),
		MinPrice: r.URL.Query().Get("min_price"),
		MaxPrice: r.URL.Query().Get("max_price"),
		InStock:  r.URL.Query().Get("in_stock"),
	}

	books, err := h.useCase.GetBooks(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	for _, b := range books {
		fmt.Fprintf(w, "ID:%d | %s | Author:%d | Price:%.2f | Stock:%d\n",
			b.ID, b.Title, b.AuthorID, b.Price, b.Stock)
	}
}

func (h *BookHandler) GetBookByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/id/")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	book, err := h.useCase.GetBookByID(idStr)
	if err != nil {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "ID: %d | %s | Price: %.2f | Stock: %d\n",
		book.ID, book.Title, book.Price, book.Stock)
}

func (h *BookHandler) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	authorID, err := strconv.Atoi(query.Get("author_id"))
	if err != nil {
		http.Error(w, "invalid author_id", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(query.Get("price"), 64)
	if err != nil {
		http.Error(w, "invalid price", http.StatusBadRequest)
		return
	}

	stock, err := strconv.Atoi(query.Get("stock"))
	if err != nil {
		http.Error(w, "invalid stock", http.StatusBadRequest)
		return
	}

	request := dto.CreateBookRequest{
		Title:    query.Get("title"),
		AuthorID: authorID,
		Price:    price,
		Stock:    stock,
	}

	msg, err := h.useCase.CreateBook(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201 Created
	fmt.Fprintln(w, msg)
}
