# Header Resolver

Extracts a value from an HTTP request header.

## Tag

```
gofast:"header:<header-name>"
```

## Example

```go
func Handler(req struct {
    Token     string `gofast:"header:Authorization"`
    RequestID string `gofast:"header:X-Request-ID"`
    Retries   int    `gofast:"header:X-Retry-Count"`
}) (*Response, error) {
    // req.Token     = "Bearer abc123"
    // req.RequestID = "req-456"
    // req.Retries   = 3  (auto-converted from string)
}
```

## Behavior

- Reads via `request.Header.Get(name)` — case-insensitive per HTTP spec
- Missing headers resolve to the zero value of the field type (empty string, 0, false)
- Automatic [type conversion](../type-conversion.md) for non-string types

## Comparison

| Framework | Code |
|-----------|------|
| **go-fast** | `Token string \`gofast:"header:Authorization"\`` |
| Gin | `token := c.GetHeader("Authorization")` |
| Fiber | `token := c.Get("Authorization")` |
