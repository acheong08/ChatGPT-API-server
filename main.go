package main

import (

	// Import local packages
	"github.com/ChatGPT-Hackers/go-server/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	//// # Add routes
	// Register new client connection
	router.GET("/client/register", handlers.ClientRegister)
	router.POST("/api/ask", handlers.ApiAsk)
}
