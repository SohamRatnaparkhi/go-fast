# DX Comparison: go-fast vs Gin vs Fiber

A side-by-side comparison of the same endpoint implemented in three frameworks.

## Example: Create Order (all five input sources)

### go-fast — 15 lines of business logic

```go
func CreateOrder(req struct {
    Body     OrderBody `gofast:"body"`
    UserID   int       `gofast:"path:user_id"`
    Currency string    `gofast:"query:currency"`
    Token    string    `gofast:"header:Authorization"`
    Session  string    `gofast:"cookie:sid"`
}) (*OrderResponse, error) {
    return &OrderResponse{
        UserID:    req.UserID,
        Item:      req.Body.Item,
        Quantity:  req.Body.Quantity,
        Price:     req.Body.Price,
        Currency:  req.Currency,
        Token:     req.Token,
        SessionID: req.Session,
    }, nil
}

h, _ := handler.Adapt(CreateOrder)
```

### Gin — 35+ lines of plumbing

```go
r.POST("/orders/:user_id", func(c *gin.Context) {
    // Path param — manual parse
    userIDStr := c.Param("user_id")
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
        return
    }

    // Header — manual read
    token := c.GetHeader("Authorization")

    // Query — manual read
    currency := c.DefaultQuery("currency", "USD")

    // Cookie — manual read + error handling
    session, err := c.Cookie("sid")
    if err != nil {
        session = ""
    }

    // Body — manual bind
    var body OrderBody
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, OrderResponse{
        UserID:    userID,
        Item:      body.Item,
        Quantity:  body.Quantity,
        Price:     body.Price,
        Currency:  currency,
        Token:     token,
        SessionID: session,
    })
})
```

### Fiber — 35+ lines of plumbing

```go
app.Post("/orders/:user_id", func(c *fiber.Ctx) error {
    userIDStr := c.Params("user_id")
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user_id"})
    }

    token := c.Get("Authorization")
    currency := c.Query("currency", "USD")
    session := c.Cookies("sid")

    var body OrderBody
    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(OrderResponse{
        UserID:    userID,
        Item:      body.Item,
        Quantity:  body.Quantity,
        Price:     body.Price,
        Currency:  currency,
        Token:     token,
        SessionID: session,
    })
})
```

## Line Count Comparison

| Task | go-fast | Gin | Fiber |
|------|---------|-----|-------|
| Extract path param + convert | 0 (tag) | 5 | 5 |
| Read header | 0 (tag) | 1 | 1 |
| Read query param | 0 (tag) | 1 | 1 |
| Read cookie | 0 (tag) | 3 | 1 |
| Parse JSON body | 0 (tag) | 4 | 4 |
| Error handling boilerplate | 0 (auto) | 6+ | 4+ |
| Write response | 0 (auto) | 1 | 1 |
| **Total plumbing** | **0** | **~20** | **~16** |

## What go-fast Eliminates

1. **Manual extraction** — No `c.Param()`, `c.Query()`, `c.GetHeader()`, `c.Cookie()` calls
2. **Manual type conversion** — No `strconv.Atoi()`, `strconv.ParseBool()`, etc.
3. **Manual error checking** — No `if err != nil { c.JSON(400, ...) }` after every extraction
4. **Manual body binding** — No `c.ShouldBindJSON(&req)` + error check
5. **Manual response writing** — No `c.JSON(200, resp)` — just return the value
6. **Framework lock-in** — Handlers are plain functions, not `func(c *gin.Context)`. Easy to test, easy to refactor.

## Testing Comparison

### go-fast — test a plain function

```go
func TestCreateOrder(t *testing.T) {
    resp, err := CreateOrder(Input{
        Body:   OrderBody{Item: "widget"},
        UserID: 42,
    })
    assert.NoError(t, err)
    assert.Equal(t, 42, resp.UserID)
}
```

### Gin — need httptest + recorder

```go
func TestCreateOrder(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("POST", "/orders/42?currency=USD", body)
    c.Params = gin.Params{{Key: "user_id", Value: "42"}}
    c.Request.Header.Set("Authorization", "Bearer tok")
    c.Request.AddCookie(&http.Cookie{Name: "sid", Value: "sess"})

    handler(c)

    assert.Equal(t, 200, w.Code)
    // ... decode response ...
}
```

go-fast handlers are just functions. Test them like functions.
