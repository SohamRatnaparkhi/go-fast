package main

import (
	"fmt"
	"net/http"

	"github.com/sohamratnaparkhi/go-fast/pkg/handler"
)

type ProfileResponse struct {
	Session string `json:"session"`
	Theme   string `json:"theme"`
}

// GetProfile — cookies are extracted automatically by tag.
func GetProfile(req struct {
	Session string `json:"cookie:session_id"`
	Theme   string `json:"cookie:theme"`
}) (*ProfileResponse, error) {
	return &ProfileResponse{
		Session: req.Session,
		Theme:   req.Theme,
	}, nil
}

func main() {
	h, err := handler.Adapt(GetProfile)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/profile", h)
	fmt.Println("go-fast server on :8080")
	fmt.Println("curl localhost:8080/profile -b 'session_id=abc123;theme=dark'")
	_ = http.ListenAndServe(":8080", nil)
}
