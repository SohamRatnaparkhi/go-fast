# Cookie Resolver

Extracts a value from an HTTP cookie.

## Tag

```
json:"cookie:<cookie-name>"
```

## Example

```go
func GetProfile(req struct {
    Session string `json:"cookie:session_id"`
    Theme   string `json:"cookie:theme"`
}) (*ProfileResponse, error) {
    // req.Session = "abc123"
    // req.Theme   = "dark"
}
```

## Behavior

- Reads via `request.Cookie(name)`
- Returns 400 if the cookie is missing (unlike header/query which return zero values)
- Automatic [type conversion](../type-conversion.md) for non-string types

## Comparison

| Framework | Code |
|-----------|------|
| **go-fast** | `Session string \`json:"cookie:session_id"\`` |
| Gin | `session, err := c.Cookie("session_id"); if err != nil { ... }` |
| Fiber | `session := c.Cookies("session_id")` |
