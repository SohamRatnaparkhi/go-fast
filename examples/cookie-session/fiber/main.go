package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type ProfileResponse struct {
	Session string `json:"session"`
	Theme   string `json:"theme"`
}

func main() {
	app := fiber.New()

	// Fiber: you manually read each cookie.
	app.Get("/profile", func(c *fiber.Ctx) error {
		session := c.Cookies("session_id")
		theme := c.Cookies("theme")

		return c.JSON(ProfileResponse{
			Session: session,
			Theme:   theme,
		})
	})

	fmt.Println("fiber server on :8080")
	log.Fatal(app.Listen(":8080"))
}
