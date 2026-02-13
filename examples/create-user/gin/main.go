package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	r := gin.Default()

	// Gin: you manually bind JSON body and read headers inside the handler.
	r.POST("/users", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		fmt.Println("Token:", token)

		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, UserResponse{
			ID:    "user_123",
			Name:  req.Name,
			Email: req.Email,
		})
	})

	fmt.Println("gin server on :8080")
	r.Run(":8080")
}
