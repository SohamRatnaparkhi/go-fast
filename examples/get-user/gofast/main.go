package main

import (
	"fmt"
	"net/http"

	"github.com/sohamratnaparkhi/go-fast/pkg/handler"
)

type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetUser — path variable :id is extracted and converted to int automatically.
func GetUser(req struct {
	ID int `json:"path:id"`
}) (*UserResponse, error) {
	return &UserResponse{
		ID:   req.ID,
		Name: fmt.Sprintf("User #%d", req.ID),
	}, nil
}

func main() {
	h, err := handler.Adapt(GetUser)
	if err != nil {
		panic(err)
	}

	// NOTE: go-fast doesn't have its own router yet.
	// In production you'd wire this through a radix-tree router that populates ctx.Params.
	// For demo purposes we use a simple wrapper that extracts the last path segment.
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		// Simulate router populating Params — a real router does this automatically.
		// This is just to show the handler works end-to-end.
		segments := splitPath(r.URL.Path)
		if len(segments) >= 2 {
			r = r.WithContext(r.Context())
			// Inject params into a custom header for now; the adapter reads ctx.Params.
			// We'll override the handler to inject params properly.
		}
		// For now, call the raw handler — path resolver will fail without ctx.Params.
		// This example is primarily to show the DX difference.
		h(w, r)
	})

	fmt.Println("go-fast server on :8080")
	fmt.Println("NOTE: Path params require a router to populate ctx.Params.")
	fmt.Println("      This example shows the handler signature DX.")
	_ = http.ListenAndServe(":8080", nil)
}

func splitPath(path string) []string {
	var parts []string
	for _, p := range split(path, '/') {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}

func split(s string, sep byte) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			if i > start {
				parts = append(parts, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		parts = append(parts, s[start:])
	}
	return parts
}
