package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := gin.Default()

	// Gin: you manually read the path param and convert it.
	r.GET("/users/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		c.JSON(http.StatusOK, UserResponse{
			ID:   id,
			Name: fmt.Sprintf("User #%d", id),
		})
	})

	fmt.Println("gin server on :8080")
	r.Run(":8080")
}
