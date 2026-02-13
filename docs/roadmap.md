# Roadmap

## Completed

- [x] **Analyzer** — Reflect-based function inspection with metadata caching
- [x] **Resolvers** — Body, Header, Query, Path, Cookie with automatic type conversion
- [x] **Adapter** — `Adapt(fn)` wiring with startup validation and per-request closure
- [x] **Type conversion** — string/bool/int*/uint*/float*/pointer support
- [x] **Error handling** — Automatic 400/500 responses from resolver and handler errors
- [x] **Examples** — Side-by-side comparisons with Gin and Fiber

## In Progress

- [ ] **Structured concurrency** (`pkg/async`) — `async.Group` with `Future[T]` for parallel work without goroutine leaks

## Planned

### Week 1: Core Engine
- [ ] **Radix tree router** — O(k) path matching with parameter extraction, populates `ctx.Params`
- [ ] **Context pooling** — `sync.Pool` for zero-alloc context reuse
- [ ] **Middleware chain** — Composable middleware with `next()` pattern
- [ ] **Validation** — Struct tag-based validation (required, min, max, pattern)
- [ ] **Dependency injection** — Constructor-based DI for services

### Week 2: Batteries
- [ ] **OpenAPI generation** — Auto-generate OpenAPI 3.0 spec from handler signatures
- [ ] **Security middleware** — CORS, CSRF, rate limiting
- [ ] **Observability** — Prometheus metrics, structured logging, request tracing
- [ ] **Advanced I/O** — Streaming responses, SSE, WebSocket support
- [ ] **Performance** — Zero-alloc hot path, benchmark suite vs Gin/Fiber/Echo

### Week 3: Polish
- [ ] **CLI tool** — `go-fast new`, `go-fast generate`, scaffolding
- [ ] **Test utilities** — `handler.Test(fn, input)` for unit testing without HTTP
- [ ] **Slice query params** — `?tag=a&tag=b` → `[]string{"a", "b"}`
- [ ] **Custom type conversion** — `encoding.TextUnmarshaler` support
- [ ] **Time parsing** — `time.Time` from string with configurable format
- [ ] **Default values** — `json:"query:page,default=1"` tag extension

## DX Improvements Proposed

### Today (implemented)
- **Declare, don't extract** — Struct tags replace manual `c.Param()`, `c.Query()`, etc.
- **Auto type conversion** — No more `strconv.Atoi()` scattered through handlers
- **Plain function handlers** — No framework context parameter, easy to test
- **Startup validation** — Bad handler signatures caught at boot, not at request time

### Coming Soon
- **Auto OpenAPI docs** — Handler signatures generate API documentation automatically
- **Zero-alloc routing** — Radix tree + context pooling for production performance
- **Parallel resolution** — Resolve independent fields concurrently via `async.Group`
- **Validation tags** — `validate:"required,min=1,max=100"` on struct fields
- **Test helpers** — `handler.Test(CreateUser, input)` returns `(output, error)` directly
