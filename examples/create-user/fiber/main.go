package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	app := fiber.New()

	// Fiber: you manually parse body and read headers inside the handler.
	app.Post("/users", func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		fmt.Println("Token:", token)

		var req CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(UserResponse{
			ID:    "user_123",
			Name:  req.Name,
			Email: req.Email,
		})
	})

	fmt.Println("fiber server on :8080")
	log.Fatal(app.Listen(":8080"))
}
