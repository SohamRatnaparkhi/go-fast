package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	app := fiber.New()

	// Fiber: you manually read the path param and convert it.
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		return c.JSON(UserResponse{
			ID:   id,
			Name: fmt.Sprintf("User #%d", id),
		})
	})

	fmt.Println("fiber server on :8080")
	log.Fatal(app.Listen(":8080"))
}
