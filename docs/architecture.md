# Architecture

## Overview

go-fast's handler system has four layers that execute in two phases:

```
STARTUP (once per handler)              RUNTIME (per request)
┌──────────┐   ┌──────────────┐        ┌─────────────────────┐
│ Analyzer │──▶│ Resolver     │──▶     │ Adapter Closure     │
│          │   │ Builder      │        │                     │
│ Inspect  │   │ Read tags,   │        │ 1. New struct       │
│ function │   │ create typed │        │ 2. Resolve fields   │
│ via      │   │ resolvers    │        │ 3. Call function    │
│ reflect  │   │              │        │ 4. Write response   │
└──────────┘   └──────────────┘        └─────────────────────┘
```

## Components

### `pkg/handler/analyzer.go` — Function Inspector

**Input:** Any `interface{}`
**Output:** `*HandlerMetadata`

Uses `reflect` to extract:
- `FuncValue` / `FuncType` — the callable and its type
- `NumInputs` / `InputTypes` — parameter count and types
- `NumOutputs` / `OutputTypes` — return value count and types
- `ReturnsError` — whether the last return implements `error`

Called once at startup. The metadata is captured in the closure.

### `pkg/handler/metadata.go` — Cached Analysis

A plain struct holding pre-computed function metadata. No methods, no logic — just data.

```go
type HandlerMetadata struct {
    FuncValue    reflect.Value
    FuncType     reflect.Type
    NumInputs    int
    NumOutputs   int
    InputTypes   []reflect.Type
    OutputTypes  []reflect.Type
    ReturnsError bool
}
```

### `pkg/handler/resolver.go` — Field Resolvers

Each resolver implements:

```go
type FieldResolver interface {
    Resolve(ctx *Context) (reflect.Value, error)
    FieldIndex() int
}
```

Seven concrete implementations:

| Resolver | Source | Key Feature |
|----------|--------|-------------|
| `BodyResolver` | `request.Body` | JSON decode into struct/pointer |
| `HeaderResolver` | `request.Header` | String + type conversion |
| `QueryResolver` | `request.URL.Query()` | String + type conversion |
| `PathVarResolver` | `ctx.Params` map | String + type conversion |
| `CookieResolver` | `request.Cookie()` | String + type conversion |
| `FormResolver` | `request.PostFormValue()` | String + type conversion |
| `FileResolver` | `request.MultipartForm.File` | `*multipart.FileHeader` |

The five string-based resolvers share `convertStringToType()` for automatic type conversion (string, bool, int*, uint*, float*, pointer). The file resolver returns `*multipart.FileHeader` directly.

### `pkg/handler/context.go` — Request Context

```go
type Context struct {
    Request  *http.Request
    Response http.ResponseWriter
    Params   map[string]string
}
```

Minimal wrapper. `Params` is populated by the router (or the adapter with an empty map as default).

### `pkg/handler/adapter.go` — The Wiring

`Adapt(fn)` orchestrates everything:

1. Calls `Analyze(fn)` → metadata
2. Validates: single struct input
3. Calls `buildResolvers(inputType)` → reads tags, creates resolvers
4. Returns a closure that uses the pre-built metadata and resolvers

The closure is a standard `http.HandlerFunc` with zero startup-time reflection.

## Data Flow

```
HTTP Request
    │
    ▼
Adapter Closure
    │
    ├── reflect.New(inputType)     → zero-value struct
    │
    ├── BodyResolver.Resolve()     → populate body field (JSON)
    ├── HeaderResolver.Resolve()   → populate header field
    ├── QueryResolver.Resolve()    → populate query field
    ├── PathVarResolver.Resolve()  → populate path field
    ├── CookieResolver.Resolve()   → populate cookie field
    ├── FormResolver.Resolve()     → populate form field
    ├── FileResolver.Resolve()     → populate file field (multipart)
    │
    ├── meta.FuncValue.Call()      → invoke handler
    │
    ├── Check error return         → 500 if error
    └── json.Encode(result)        → 200 JSON response
```

## Design Decisions

1. **Struct tags over conventions** — Explicit is better than implicit. `gofast:"header:Authorization"` is unambiguous.
2. **One struct input** — Forces grouping of all inputs. Makes the handler self-documenting.
3. **Startup validation** — `Adapt()` returns errors for bad signatures. No runtime surprises.
4. **Body resolved first** — Body consumes `request.Body` (a reader), so it must run before anything else that might need it.
5. **Standard `http.HandlerFunc`** — The output of `Adapt()` works with any Go HTTP router or middleware.
6. **Body vs form/file exclusivity** — JSON body and form/file resolvers both consume the request body through different parsers. `Adapt()` rejects structs that mix them.
