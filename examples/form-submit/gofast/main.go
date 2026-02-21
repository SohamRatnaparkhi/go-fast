package main

import (
	"fmt"
	"net/http"

	"github.com/sohamratnaparkhi/go-fast/pkg/handler"
)

type ContactResponse struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// SubmitContact handles a URL-encoded or multipart form submission.
//
// Each form field is declared as a struct field with a form: tag.
// Type conversion works the same as query/header resolvers.
func SubmitContact(req struct {
	Name    string `gofast:"form:name"`
	Email   string `gofast:"form:email"`
	Message string `gofast:"form:message"`
}) (*ContactResponse, error) {
	return &ContactResponse{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
	}, nil
}

func main() {
	h, err := handler.Adapt(SubmitContact)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/contact", h)
	fmt.Println("go-fast server on :8080")
	fmt.Println(`
Form submission — zero boilerplate:

  curl -X POST localhost:8080/contact \
    -d 'name=Alice&email=alice@example.com&message=Hello!'`)
	_ = http.ListenAndServe(":8080", nil)
}
