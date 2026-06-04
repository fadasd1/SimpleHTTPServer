package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// GET /books/{id}
	resp, err := http.Get("http://api:8080/books/1")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	// POST /books — JSON body instead of query params
	createPayload, _ := json.Marshal(map[string]any{
		"title":     "NewBook",
		"author_id": 1,
		"price":     100,
		"stock":     25,
	})
	resp, err = http.Post("http://api:8080/books", "application/json", bytes.NewBuffer(createPayload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ = io.ReadAll(resp.Body)
	fmt.Println(string(body))

	// POST /books/list — JSON body instead of query params
	listPayload, _ := json.Marshal(map[string]any{
		"page":     "1",
		"limit":    "10",
		"in_stock": "true",
	})
	resp, err = http.Post("http://api:8080/books/list", "application/json", bytes.NewBuffer(listPayload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ = io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
