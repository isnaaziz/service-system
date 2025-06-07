package controllers

import (
	"service_system/services"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	userService := services.UserService{DB: db}
	users, err := userService.GetActiveUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": "Error retrieving users"})
		return
	}

	c.JSON(200, users)
}

func SoftDeleteUser(c *gin.Context) {
	userID := c.Param("id")

	userService := services.UserService{DB: db}
	if err := userService.SoftDeleteUser(userID); err != nil {
		c.JSON(500, gin.H{"error": "Error soft deleting user"})
		return
	}

	c.JSON(200, gin.H{"message": "User soft deleted successfully"})
}
