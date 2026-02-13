# Resolvers

Resolvers are the core of go-fast's DX improvement. They automatically extract and convert request data into your handler's input struct fields based on struct tags.

## How It Works

When you call `handler.Adapt(fn)`, go-fast:

1. Inspects the input struct's fields and their `json` tags
2. Builds a typed resolver for each tagged field
3. At request time, each resolver extracts its value from the correct source and converts it to the declared Go type

## Tag Format

```
json:"<source>:<name>"
```

| Source | Tag | Reads From |
|--------|-----|------------|
| Body | `json:"body"` | `request.Body` (JSON) |
| Header | `json:"header:<name>"` | `request.Header.Get(name)` |
| Query | `json:"query:<name>"` | `request.URL.Query().Get(name)` |
| Path | `json:"path:<name>"` | `ctx.Params[name]` |
| Cookie | `json:"cookie:<name>"` | `request.Cookie(name)` |

## Example: All Five in One Handler

```go
func CreateOrder(req struct {
    Body     OrderBody `json:"body"`
    UserID   int       `json:"path:user_id"`
    Currency string    `json:"query:currency"`
    Token    string    `json:"header:Authorization"`
    Session  string    `json:"cookie:sid"`
}) (*OrderResponse, error) {
    // req.Body     — parsed from JSON body
    // req.UserID   — extracted from URL path, converted to int
    // req.Currency — extracted from ?currency=...
    // req.Token    — extracted from Authorization header
    // req.Session  — extracted from sid cookie
}
```

## Rules

- Fields must be **exported** (uppercase first letter)
- Only one `json:"body"` field is allowed per struct
- Tag names cannot be empty (e.g., `json:"header:"` is invalid)
- Untagged or `json:"-"` fields are skipped
- String-based resolvers (header, query, path, cookie) support automatic [type conversion](../type-conversion.md)

## Detailed Docs

- [Body Resolver](./body.md)
- [Header Resolver](./header.md)
- [Query Resolver](./query.md)
- [Path Resolver](./path.md)
- [Cookie Resolver](./cookie.md)
