package main

import (
	"fmt"
	"net/http"

	"github.com/sohamratnaparkhi/go-fast/pkg/handler"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUser — the handler is a plain function.
// go-fast inspects struct tags at startup, wires resolvers, and calls this automatically.
func CreateUser(req struct {
	Body  CreateUserRequest `json:"body"`
	Token string            `json:"header:Authorization"`
}) (*UserResponse, error) {
	fmt.Println("Token:", req.Token)
	return &UserResponse{
		ID:    "user_123",
		Name:  req.Body.Name,
		Email: req.Body.Email,
	}, nil
}

func main() {
	h, err := handler.Adapt(CreateUser)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/users", h)
	fmt.Println("go-fast server on :8080")
	fmt.Println("curl -X POST localhost:8080/users -H 'Authorization: Bearer tok' -d '{\"name\":\"John\",\"email\":\"j@test.com\"}'")
	_ = http.ListenAndServe(":8080", nil)
}
