package controllers

import (
	"service_system/services"

	"github.com/gin-gonic/gin"
)

// GetUsers godoc
// @Summary Get all users
// @Description Get all active users
// @Tags User
// @Produce  json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(c *gin.Context) {
	userService := services.UserService{DB: db}
	users, err := userService.GetActiveUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": "Error retrieving users"})
		return
	}

	c.JSON(200, users)
}

// SoftDeleteUser godoc
// @Summary Soft delete user
// @Description Soft delete user by ID
// @Tags User
// @Param   id  path  int  true  "User ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [delete]
func SoftDeleteUser(c *gin.Context) {
	userID := c.Param("id")

	userService := services.UserService{DB: db}
	if err := userService.SoftDeleteUser(userID); err != nil {
		c.JSON(500, gin.H{"error": "Error soft deleting user"})
		return
	}

	c.JSON(200, gin.H{"message": "User soft deleted successfully"})
}
