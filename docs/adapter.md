# Adapter

The adapter is the glue that connects your handler function to the HTTP server. It's the single function you call: `handler.Adapt(fn)`.

## API

```go
func Adapt(fn interface{}) (http.HandlerFunc, error)
```

Takes any function with a single struct input, returns a standard `http.HandlerFunc`.

## What It Does

### At Startup (once)

1. **Analyze** — Inspects the function via reflection: input types, output types, whether it returns an error
2. **Build resolvers** — Reads struct tags on the input type, creates a typed resolver for each field
3. **Validate** — Rejects invalid signatures (non-function, multiple inputs, non-struct input, duplicate body tags, empty tag names, unexported tagged fields)

### Per Request (closure)

1. Create a new zero-value instance of the input struct
2. Resolve the body field first (if present)
3. Resolve all other fields (header, query, path, cookie)
4. Call the handler function with the populated struct
5. If the function returns an error, write a 500 JSON error response
6. If the function returns a value, write it as JSON with 200
7. If the function returns nothing (no non-error outputs), write 204 No Content

## Handler Signatures

### Value + Error (most common)

```go
func Handler(req Input) (*Output, error)
```

### Value Only (no error possible)

```go
func Handler(req Input) *Output
```

### Error Only (side-effect handler)

```go
func Handler(req Input) error
```

### No Return (fire-and-forget)

```go
func Handler(req Input)
```

Returns 204 No Content automatically.

## Error Handling

| Error Source | HTTP Status | When |
|-------------|-------------|------|
| Malformed body JSON | 400 | Body resolver fails to decode |
| Missing path variable | 400 | Path param not in `ctx.Params` |
| Missing cookie | 400 | Cookie not present in request |
| Type conversion failure | 400 | e.g., `"abc"` for an `int` field |
| Handler returns error | 500 | Last return value is non-nil error |
| Response encoding failure | 500 | JSON marshal of return value fails |

All errors are returned as JSON:

```json
{"error": "decode body: unexpected EOF"}
```

## Startup Validation Errors

`Adapt()` returns an error (not a panic) for these cases:

- Argument is not a function
- Function has != 1 input parameter
- Input parameter is not a struct
- Tagged field is unexported
- Multiple `json:"body"` fields
- Empty tag name (e.g., `json:"header:"`)

This means invalid handlers are caught at server startup, not at request time.
