# Query Resolver

Extracts a value from URL query parameters.

## Tag

```
gofast:"query:<param-name>"
```

## Example

```go
func SearchUsers(req struct {
    Query  string `gofast:"query:q"`
    Page   int    `gofast:"query:page"`
    Active bool   `gofast:"query:active"`
}) (*SearchResult, error) {
    // GET /search?q=john&page=2&active=true
    // req.Query  = "john"
    // req.Page   = 2      (auto-converted)
    // req.Active = true   (auto-converted)
}
```

## Behavior

- Reads via `request.URL.Query().Get(name)`
- Missing params resolve to the zero value (empty string, 0, false)
- Automatic [type conversion](../type-conversion.md) for int, uint, float, bool types
- Only reads the first value for a given key (no multi-value support yet)

## Comparison

| Framework | Code |
|-----------|------|
| **go-fast** | `Page int \`gofast:"query:page"\`` |
| Gin | `pageStr := c.Query("page"); page, _ := strconv.Atoi(pageStr)` |
| Fiber | `pageStr := c.Query("page"); page, _ := strconv.Atoi(pageStr)` |

go-fast eliminates the manual `strconv` call. Declare the type, get the type.
