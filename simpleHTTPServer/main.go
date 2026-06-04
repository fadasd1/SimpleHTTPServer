package main

import (
	"context"
	"log"
	"net/http"
	"simpleHTTPServer/controller"
	"simpleHTTPServer/repository"
	"simpleHTTPServer/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://fadasd:1337@dummyDB:5432/dummyDB?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}
	defer pool.Close()

	repoConnection := repository.NewBookRepository(pool)

	uc := usecase.NewBookUseCase(repoConnection)
	h := controller.NewBookHandler(uc)

	http.HandleFunc("GET /books/{id}", h.GetBookByIDHandler)
	http.HandleFunc("POST /books/list", h.GetBooksHandler)
	http.HandleFunc("POST /books", h.CreateBookHandler)

	log.Println("API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
