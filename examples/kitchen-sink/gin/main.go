package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	r := gin.Default()

	// Gin: every source requires manual extraction and type conversion.
	r.POST("/orders/:user_id", func(c *gin.Context) {
		// Path param — manual parse
		userIDStr := c.Param("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}

		// Header — manual read
		token := c.GetHeader("Authorization")

		// Query — manual read
		currency := c.DefaultQuery("currency", "USD")

		// Cookie — manual read + error handling
		session, err := c.Cookie("sid")
		if err != nil {
			session = ""
		}

		// Body — manual bind
		var body OrderBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, OrderResponse{
			UserID:    userID,
			Item:      body.Item,
			Quantity:  body.Quantity,
			Price:     body.Price,
			Currency:  currency,
			Token:     token,
			SessionID: session,
		})
	})

	fmt.Println("gin server on :8080")
	r.Run(":8080")
}
