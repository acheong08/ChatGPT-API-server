package handlers

import (
	"os"

	"github.com/ChatGPT-Hackers/ChatGPT-API-server/utils"
	_ "github.com/ChatGPT-Hackers/ChatGPT-API-server/utils"
	"github.com/gin-gonic/gin"
)

type Request struct {
	AdminKey string `json:"admin_key"`
	UserID   string `json:"user_id"`
}

func Admin_userAdd(c *gin.Context) {
	// Get admin key from request body
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Check if admin key is valid
	if !utils.VerifyAdminKey(request.AdminKey) {
		c.JSON(401, gin.H{
			"error": "Invalid admin key",
		})
		return
	}

	// Generate user_id and token
	user_id := utils.GenerateId()
	token := utils.GenerateId()

	// Insert user_id and token into database
	err := utils.DatabaseInsert(user_id, token)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to insert user_id and token into database",
		})
		return
	}

	// Return user_id and token
	c.JSON(200, gin.H{
		"user_id": user_id,
		"token":   token,
	})
}

// POST request to delete a user
func Admin_userDel(c *gin.Context) {
	// Get admin key from request body
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Check if admin key is valid
	if !utils.VerifyAdminKey(request.AdminKey) {
		c.JSON(401, gin.H{
			"error": "Invalid admin key",
		})
		return
	}

	// Delete user from database
	err := utils.DatabaseDelete(request.UserID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to delete user from database",
		})
		return
	}

	// Return success
	c.JSON(200, gin.H{
		"message": "User deleted",
	})
}

func Admin_usersGet(c *gin.Context) {
	// Get admin key from GET parameter
	AdminKey := c.Query("admin_key")

	// Check if admin key is valid
	if !utils.VerifyAdminKey(AdminKey) {
		c.JSON(401, gin.H{
			"error":   "Invalid admin key",
			"key":     AdminKey,
			"correct": os.Args[2],
		})
		return
	}

	// Get users from database
	users, err := utils.DatabaseSelectAll()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to get users from database",
			"error":   err.Error(),
		})
		return
	}

	// Return users
	c.JSON(200, gin.H{
		"users": users,
	})
}
