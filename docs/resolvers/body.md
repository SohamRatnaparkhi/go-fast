# Body Resolver

Parses the JSON request body into a struct field.

## Tag

```
gofast:"body"
```

## Supported Field Types

- Any struct type (`CreateUserRequest`)
- Pointer to struct (`*CreateUserRequest`)

## Example

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func CreateUser(req struct {
    Body CreateUserRequest `gofast:"body"`
}) (*UserResponse, error) {
    // req.Body.Name and req.Body.Email are populated from JSON
}
```

## Pointer Body

```go
func CreateUser(req struct {
    Body *CreateUserRequest `gofast:"body"`
}) (*UserResponse, error) {
    // req.Body is a *CreateUserRequest — nil check if needed
}
```

## Behavior

- Uses `json.NewDecoder(request.Body).Decode(...)` under the hood
- Returns 400 if the body is malformed JSON
- Only one body field is allowed per input struct — a second `gofast:"body"` tag causes a startup error
- The body struct uses standard `json` tags for field mapping (`json:"name"`, `json:"email"`, etc.)

## Comparison

| Framework | Code |
|-----------|------|
| **go-fast** | `Body CreateUserRequest \`gofast:"body"\`` |
| Gin | `var req CreateUserRequest; c.ShouldBindJSON(&req)` |
| Fiber | `var req CreateUserRequest; c.BodyParser(&req)` |
