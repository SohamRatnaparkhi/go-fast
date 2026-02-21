# Overview

## What is go-fast?

go-fast is a Go HTTP framework that eliminates handler boilerplate. You write plain functions with a struct input — go-fast reads struct tags, resolves parameters from the request automatically, calls your function, and writes the response.

**The problem with every Go framework today:**

```go
// Gin, Echo, Fiber, Chi — they all look like this
func CreateUser(c *gin.Context) {
    token := c.GetHeader("Authorization")       // manual
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil { // manual
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    user, err := createUser(req)
    if err != nil {                               // manual
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, user)                             // manual
}
```

Every handler repeats the same pattern: extract, parse, convert, validate, call, check error, write response. It's tedious, error-prone, and obscures business logic.

**go-fast:**

```go
func CreateUser(req struct {
    Body  CreateUserRequest `gofast:"body"`
    Token string            `gofast:"header:Authorization"`
}) (*UserResponse, error) {
    fmt.Println("Token:", req.Token)
    return &UserResponse{ID: "123", Name: req.Body.Name}, nil
}

h, _ := handler.Adapt(CreateUser)
http.HandleFunc("/users", h)
```

That's it. No manual parsing, no manual error handling, no manual response writing. Declare what you need, get it.

## Core Principles

1. **Declare, don't extract** — Struct tags describe where data comes from. The framework does the rest.
2. **Analyze once, run fast** — Reflection happens at startup. The per-request closure is pre-wired.
3. **Plain functions** — Handlers are regular Go functions. No framework-specific context parameter. Easy to test, easy to read.
4. **Type-safe** — Query params, path vars, headers, and cookies are automatically converted to the declared Go type (int, bool, float64, etc.).
5. **Zero magic at runtime** — Everything is explicit via struct tags. No hidden conventions.

## Features (Current)

| Feature | Tag | Example |
|---------|-----|---------|
| JSON body | `gofast:"body"` | `Body CreateUserRequest \`gofast:"body"\`` |
| Header | `gofast:"header:<name>"` | `Token string \`gofast:"header:Authorization"\`` |
| Query param | `gofast:"query:<name>"` | `Page int \`gofast:"query:page"\`` |
| Path variable | `gofast:"path:<name>"` | `ID int \`gofast:"path:id"\`` |
| Cookie | `gofast:"cookie:<name>"` | `Session string \`gofast:"cookie:session_id"\`` |
| Auto type conversion | — | string, int*, uint*, float*, bool, pointer types |
| Error return handling | — | Return `error` as last value → automatic 500 |
| No-return handlers | — | Return nothing → automatic 204 No Content |
