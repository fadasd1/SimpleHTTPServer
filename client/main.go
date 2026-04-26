package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	url := "http://api:8080/books/id/1"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))

	url = "http://api:8080/books/create?title=NewBook&author_id=1&price=100&stock=25"

	resp, err = http.Get(url)
	if err != nil {
		panic(err)
	}
	body, _ = io.ReadAll(resp.Body)

	fmt.Println(string(body))

	url = "http://api:8080/books/get?page=1&limit=10&in_stock=true"

	resp, err = http.Get(url)
	if err != nil {
		panic(err)
	}
	body, _ = io.ReadAll(resp.Body)

	fmt.Println(string(body))
}
