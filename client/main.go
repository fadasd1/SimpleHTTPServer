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
}
