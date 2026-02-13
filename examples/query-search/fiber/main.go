package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SearchResult struct {
	Query  string `json:"query"`
	Page   int    `json:"page"`
	Active bool   `json:"active"`
}

func main() {
	app := fiber.New()

	// Fiber: you manually read each query param and convert types yourself.
	app.Get("/search", func(c *fiber.Ctx) error {
		q := c.Query("q")
		pageStr := c.Query("page", "0")
		activeStr := c.Query("active", "false")

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid page"})
		}

		active, err := strconv.ParseBool(activeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid active"})
		}

		return c.JSON(SearchResult{
			Query:  q,
			Page:   page,
			Active: active,
		})
	})

	fmt.Println("fiber server on :8080")
	log.Fatal(app.Listen(":8080"))
}
