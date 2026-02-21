# Path Variable Resolver

Extracts a value from URL path parameters (e.g., `/users/:id`).

## Tag

```
gofast:"path:<param-name>"
```

## Example

```go
func GetUser(req struct {
    ID int `gofast:"path:id"`
}) (*UserResponse, error) {
    // GET /users/42
    // req.ID = 42 (auto-converted from string)
}
```

## Behavior

- Reads from `ctx.Params[name]` — a `map[string]string` populated by the router
- Returns 400 if the param is missing from the map
- Automatic [type conversion](../type-conversion.md) for non-string types
- **Requires a router** that populates `ctx.Params` before the handler runs (go-fast's radix tree router is on the roadmap)

## Comparison

| Framework | Code |
|-----------|------|
| **go-fast** | `ID int \`gofast:"path:id"\`` |
| Gin | `idStr := c.Param("id"); id, _ := strconv.Atoi(idStr)` |
| Fiber | `idStr := c.Params("id"); id, _ := strconv.Atoi(idStr)` |
