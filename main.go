package main

import (

	// Import local packages

	"net/http"

	"os"

	"github.com/ChatGPT-Hackers/ChatGPT-API-server/handlers"
	"github.com/ChatGPT-Hackers/ChatGPT-API-server/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// get arg server port and admin key
	if len(os.Args) < 3 {
		println("Usage: ./ChatGPT-API-server <port> <admin key>")
		return
	}
	println(os.Args[1], os.Args[2])

	// Make database
	err := utils.DatabaseCreate()
	if err != nil {
		println("Failed to create database:", err.Error())
		return
	}

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
	router.POST("/admin/users/add", handlers.Admin_userAdd)
	router.POST("/admin/users/delete", handlers.Admin_userDel)
	router.GET("/admin/users", handlers.Admin_usersGet)

	// Add a health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Start server
	router.Run(":" + os.Args[1])
}
