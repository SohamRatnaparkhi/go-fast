package main

import (
	"fmt"
	"net/http"

	"github.com/sohamratnaparkhi/go-fast/pkg/handler"
)

type OrderBody struct {
	Item     string  `json:"item"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type OrderResponse struct {
	UserID    int     `json:"user_id"`
	Item      string  `json:"item"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
	Token     string  `json:"token"`
	SessionID string  `json:"session_id"`
}

// CreateOrder uses ALL five resolver types in a single handler:
//   - body:           JSON request body
//   - path:user_id    path variable (requires router)
//   - query:currency  query parameter
//   - header:Authorization  request header
//   - cookie:sid      cookie value
func CreateOrder(req struct {
	Body    OrderBody `gofast:"body"`
	UserID  int       `gofast:"path:user_id"`
	Currency string   `gofast:"query:currency"`
	Token   string    `gofast:"header:Authorization"`
	Session string    `gofast:"cookie:sid"`
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

func main() {
	h, err := handler.Adapt(CreateOrder)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/orders", h)
	fmt.Println("go-fast server on :8080")
	fmt.Println(`
All five resolvers in one handler — zero boilerplate:

  curl -X POST 'localhost:8080/orders?currency=USD' \
    -H 'Authorization: Bearer tok' \
    -b 'sid=sess-abc' \
    -d '{"item":"widget","quantity":3,"price":9.99}'

NOTE: path:user_id requires a router to populate ctx.Params.
      This example shows the DX; path resolution works end-to-end
      once a router is wired.`)
	_ = http.ListenAndServe(":8080", nil)
}
