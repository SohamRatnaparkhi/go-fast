# File Resolver

Extracts an uploaded file from a `multipart/form-data` request.

## Tag

```
gofast:"file:<field-name>"
```

## Field Type

The field **must** be `*multipart.FileHeader`. This is validated at startup by `Adapt()`.

```go
import "mime/multipart"
```

## Example

```go
func UploadAvatar(req struct {
    Avatar *multipart.FileHeader `gofast:"file:avatar"`
}) (*Response, error) {
    // req.Avatar.Filename = "photo.png"
    // req.Avatar.Size     = 4096

    f, err := req.Avatar.Open()
    if err != nil {
        return nil, err
    }
    defer f.Close()
    // read file content from f
}
```

## Behavior

- Calls `request.ParseMultipartForm(32MB)` then reads from `request.MultipartForm.File`
- Returns 400 if the file is missing (required — like path and cookie)
- Returns 400 if the content type is not `multipart/form-data`
- Does **not** open the file eagerly — the `*FileHeader` is returned directly so you control when to read
- **Cannot be combined with `gofast:"body"`** — both consume the request body
- **Can** be combined with `gofast:"form:..."` for mixed field+file multipart uploads

## Startup Validation

`Adapt()` rejects the handler at startup if:
- The file field type is not `*multipart.FileHeader`
- The struct mixes `gofast:"body"` with `gofast:"file:..."`

## Comparison

| Framework | Code |
|-----------|------|
| **go-fast** | `Avatar *multipart.FileHeader \`gofast:"file:avatar"\`` |
| Gin | `file, err := c.FormFile("avatar"); if err != nil { ... }` |
| Fiber | `file, err := c.FormFile("avatar"); if err != nil { ... }` |
