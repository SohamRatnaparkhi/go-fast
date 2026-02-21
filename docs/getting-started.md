# Getting Started

## Install

```bash
go get github.com/sohamratnaparkhi/go-fast
```

## Write Your First Handler

```go
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
    Body  CreateUserRequest `gofast:"body"`
    Token string            `gofast:"header:Authorization"`
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
    fmt.Println("Server on :8080")
    http.ListenAndServe(":8080", nil)
}
```

## Test It

```bash
go run main.go

# In another terminal:
curl -X POST localhost:8080/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer my-token" \
  -d '{"name":"John","email":"john@test.com"}'
```

**Response:**
```json
{"id":"user_123","name":"John","email":"john@test.com"}
```

**Server log:**
```
Token: Bearer my-token
```

## How It Works

1. **`handler.Adapt(CreateUser)`** inspects the function signature at startup using reflection.
2. It reads struct tags on the input parameter and builds a resolver for each field.
3. It returns an `http.HandlerFunc` closure that, per request:
   - Creates a new instance of the input struct
   - Resolves each field from the appropriate request source (body, header, query, path, cookie)
   - Calls your function with the populated struct
   - Writes the return value as JSON (or the error as a 500)

Reflection only runs once. The per-request path is a pre-built closure with no reflection overhead beyond `reflect.Value.Call`.

## Next Steps

- [Resolver Tags](./resolvers/README.md) — All five resolver types
- [Type Conversion](./type-conversion.md) — Supported types
- [DX Comparison](./dx-comparison.md) — See the difference vs Gin/Fiber
