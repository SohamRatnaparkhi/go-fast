# go-fast Examples

Side-by-side comparison of **go-fast** vs **Gin** vs **Fiber** for the same endpoints.

Each example is its own Go module so framework-specific dependencies (gin, fiber) never leak into the library.

## Structure

```
examples/
├── create-user/       # Body + Header resolver
│   ├── gofast/
│   ├── gin/
│   └── fiber/
├── query-search/      # Query param resolver
│   ├── gofast/
│   ├── gin/
│   └── fiber/
├── get-user/          # Path variable resolver
│   ├── gofast/
│   ├── gin/
│   └── fiber/
├── cookie-session/    # Cookie resolver
│   ├── gofast/
│   ├── gin/
│   └── fiber/
└── kitchen-sink/      # ALL resolvers in one handler
    ├── gofast/
    ├── gin/
    └── fiber/
```

## Running

### go-fast examples (no external deps)

```bash
cd examples/create-user/gofast && go run main.go
cd examples/query-search/gofast && go run main.go
cd examples/cookie-session/gofast && go run main.go
cd examples/kitchen-sink/gofast && go run main.go
```

### Gin / Fiber examples (requires `go mod tidy` first)

```bash
cd examples/create-user/gin && go mod tidy && go run main.go
cd examples/create-user/fiber && go mod tidy && go run main.go
```

## The DX Difference

### go-fast — declare what you need, get it automatically

```go
func CreateOrder(req struct {
    Body     OrderBody `gofast:"body"`
    UserID   int       `gofast:"path:user_id"`
    Currency string    `gofast:"query:currency"`
    Token    string    `gofast:"header:Authorization"`
    Session  string    `gofast:"cookie:sid"`
}) (*OrderResponse, error) {
    // Just use req.Body, req.UserID, req.Currency, etc.
}
```

### Gin — manually extract and convert everything

```go
r.POST("/orders/:user_id", func(c *gin.Context) {
    userIDStr := c.Param("user_id")
    userID, err := strconv.Atoi(userIDStr)  // manual conversion
    if err != nil { ... }

    token := c.GetHeader("Authorization")   // manual read
    currency := c.DefaultQuery("currency", "USD")
    session, _ := c.Cookie("sid")           // manual read + error

    var body OrderBody
    if err := c.ShouldBindJSON(&body); err != nil { ... }  // manual bind

    c.JSON(200, OrderResponse{ ... })
})
```

### Fiber — same manual work, different API

```go
app.Post("/orders/:user_id", func(c *fiber.Ctx) error {
    userIDStr := c.Params("user_id")
    userID, _ := strconv.Atoi(userIDStr)
    token := c.Get("Authorization")
    currency := c.Query("currency", "USD")
    session := c.Cookies("sid")

    var body OrderBody
    c.BodyParser(&body)

    return c.JSON(OrderResponse{ ... })
})
```

**go-fast eliminates the boilerplate.** You declare your inputs via struct tags, and the framework resolves them automatically at zero runtime overhead (analysis happens once at startup).
