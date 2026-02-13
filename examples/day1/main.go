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

func CreateUser(req struct {
	Body  CreateUserRequest `json:"body"`
	Token string            `json:"header:Authorization"`
}) (*UserResponse, error) {
	fmt.Println("Token received:", req.Token)
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
	fmt.Println("Server on :8080")
	_ = http.ListenAndServe(":8080", nil)
}
