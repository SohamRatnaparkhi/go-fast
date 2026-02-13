package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SearchResult struct {
	Query  string `json:"query"`
	Page   int    `json:"page"`
	Active bool   `json:"active"`
}

func main() {
	r := gin.Default()

	// Gin: you manually read each query param and convert types yourself.
	r.GET("/search", func(c *gin.Context) {
		q := c.Query("q")
		pageStr := c.DefaultQuery("page", "0")
		activeStr := c.DefaultQuery("active", "false")

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
			return
		}

		active, err := strconv.ParseBool(activeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid active"})
			return
		}

		c.JSON(http.StatusOK, SearchResult{
			Query:  q,
			Page:   page,
			Active: active,
		})
	})

	fmt.Println("gin server on :8080")
	r.Run(":8080")
}
