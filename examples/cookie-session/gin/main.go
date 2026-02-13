package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileResponse struct {
	Session string `json:"session"`
	Theme   string `json:"theme"`
}

func main() {
	r := gin.Default()

	// Gin: you manually read each cookie and handle the error.
	r.GET("/profile", func(c *gin.Context) {
		session, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing session_id cookie"})
			return
		}

		theme, err := c.Cookie("theme")
		if err != nil {
			theme = ""
		}

		c.JSON(http.StatusOK, ProfileResponse{
			Session: session,
			Theme:   theme,
		})
	})

	fmt.Println("gin server on :8080")
	r.Run(":8080")
}
