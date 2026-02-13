package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
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

func main() {
	app := fiber.New()

	// Fiber: every source requires manual extraction and type conversion.
	app.Post("/orders/:user_id", func(c *fiber.Ctx) error {
		// Path param — manual parse
		userIDStr := c.Params("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user_id"})
		}

		// Header — manual read
		token := c.Get("Authorization")

		// Query — manual read
		currency := c.Query("currency", "USD")

		// Cookie — manual read
		session := c.Cookies("sid")

		// Body — manual parse
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

	fmt.Println("fiber server on :8080")
	log.Fatal(app.Listen(":8080"))
}
