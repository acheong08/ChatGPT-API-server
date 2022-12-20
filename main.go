package main

import (

	// Import local packages

	"github.com/ChatGPT-Hackers/go-server/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//// # Headers
	// Allow CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	})
	//// # Add routes
	// Register new client connection
	router.GET("/client/register", handlers.Client_register)
	router.POST("/api/ask", handlers.API_ask)
	router.GET("/api/connections", handlers.API_getConnections)

	// Start server
	router.Run(":8080")
}
