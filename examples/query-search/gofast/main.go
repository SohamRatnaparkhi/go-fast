package main

import (
	"fmt"
	"net/http"

	"github.com/sohamratnaparkhi/go-fast/pkg/handler"
)

type SearchResult struct {
	Query  string `json:"query"`
	Page   int    `json:"page"`
	Active bool   `json:"active"`
}

// SearchUsers — query params are extracted automatically from ?q=...&page=...&active=...
func SearchUsers(req struct {
	Query  string `json:"query:q"`
	Page   int    `json:"query:page"`
	Active bool   `json:"query:active"`
}) (*SearchResult, error) {
	return &SearchResult{
		Query:  req.Query,
		Page:   req.Page,
		Active: req.Active,
	}, nil
}

func main() {
	h, err := handler.Adapt(SearchUsers)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/search", h)
	fmt.Println("go-fast server on :8080")
	fmt.Println("curl 'localhost:8080/search?q=john&page=2&active=true'")
	_ = http.ListenAndServe(":8080", nil)
}
