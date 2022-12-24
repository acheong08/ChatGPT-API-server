package main

import (

	// Import local packages

	"os"

	"github.com/ChatGPT-Hackers/ChatGPT-API-server/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// get arg server port and secret key
	if len(os.Args) < 3 {
		println("Usage: ./ChatGPT-API-server <port> <secret key>")
		return
	}
	println(os.Args[1], os.Args[2])
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
	router.Run(":" + os.Args[1])
}
