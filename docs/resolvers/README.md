# Resolvers

Resolvers are the core of go-fast's DX improvement. They automatically extract and convert request data into your handler's input struct fields based on struct tags.

## How It Works

When you call `handler.Adapt(fn)`, go-fast:

1. Inspects the input struct's fields and their `gofast` tags
2. Builds a typed resolver for each tagged field
3. At request time, each resolver extracts its value from the correct source and converts it to the declared Go type

## Tag Format

```
gofast:"<source>:<name>"
```

| Source | Tag | Reads From |
|--------|-----|------------|
| Body | `gofast:"body"` | `request.Body` (JSON) |
| Header | `gofast:"header:<name>"` | `request.Header.Get(name)` |
| Query | `gofast:"query:<name>"` | `request.URL.Query().Get(name)` |
| Path | `gofast:"path:<name>"` | `ctx.Params[name]` |
| Cookie | `gofast:"cookie:<name>"` | `request.Cookie(name)` |
| Form | `gofast:"form:<name>"` | `request.PostFormValue(name)` |
| File | `gofast:"file:<name>"` | `request.MultipartForm.File[name]` |

## Example: All Seven in One Handler

```go
func CreateOrder(req struct {
    Body     OrderBody `gofast:"body"`
    UserID   int       `gofast:"path:user_id"`
    Currency string    `gofast:"query:currency"`
    Token    string    `gofast:"header:Authorization"`
    Session  string    `gofast:"cookie:sid"`
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
- Only one `gofast:"body"` field is allowed per struct
- Tag names cannot be empty (e.g., `gofast:"header:"` is invalid)
- Untagged or `gofast:"-"` fields are skipped
- String-based resolvers (header, query, path, cookie, form) support automatic [type conversion](../type-conversion.md)
- File fields must be `*multipart.FileHeader`
- `gofast:"body"` cannot be combined with `gofast:"form:..."` or `gofast:"file:..."` (both consume the request body)

## Detailed Docs

- [Body Resolver](./body.md)
- [Header Resolver](./header.md)
- [Query Resolver](./query.md)
- [Path Resolver](./path.md)
- [Cookie Resolver](./cookie.md)
- [Form Resolver](./form.md)
- [File Resolver](./file.md)
