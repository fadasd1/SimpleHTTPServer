package main

import (
	"database/sql"
	"log"
	"net/http"
	"simpleHTTPServer/controller"
	"simpleHTTPServer/repository"
	"simpleHTTPServer/usecase"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	var db *sql.DB
	var err error
	connStr := "postgres://fadasd:1337@dummyDB:5432/dummyDB?sslmode=disable"

	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repoConnection := repository.NewBookRepository(db)
	uc := usecase.NewBookUseCase(repoConnection)
	h := controller.NewBookHandler(uc)

	http.HandleFunc("/books/id/", h.GetBookByIDHandler)
	http.HandleFunc("/books/get", h.GetBooksHandler)
	http.HandleFunc("books/create", h.CreateBookHandler)

	log.Println("API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
