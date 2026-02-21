# Form Resolver

Extracts a value from the POST body of a form submission.

## Tag

```
json:"form:<field-name>"
```

## Example

```go
func SubmitForm(req struct {
    Name  string `json:"form:name"`
    Email string `json:"form:email"`
    Age   int    `json:"form:age"`
}) (*Response, error) {
    // req.Name  = "alice"
    // req.Email = "alice@example.com"
    // req.Age   = 30
}
```

## Behavior

- Reads via `request.PostFormValue(name)` (POST body only, not URL query)
- Works with both `application/x-www-form-urlencoded` and `multipart/form-data`
- Missing fields return zero values (like header/query — not an error)
- Automatic [type conversion](../type-conversion.md) for non-string types
- **Cannot be combined with `json:"body"`** — both consume the request body

## Combined with File Resolver

Form and file resolvers work together for multipart uploads:

```go
func UploadDocument(req struct {
    Title    string                `json:"form:title"`
    Document *multipart.FileHeader `json:"file:document"`
}) (*Response, error) {
    // req.Title    = "My Report"
    // req.Document = uploaded file header
}
```

## Comparison

| Framework | Code |
|-----------|------|
| **go-fast** | `Name string \`json:"form:name"\`` |
| Gin | `name := c.PostForm("name")` |
| Fiber | `name := c.FormValue("name")` |
