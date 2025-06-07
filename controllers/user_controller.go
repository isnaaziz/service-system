package controllers

import (
	"service_system/middleware"
	"service_system/services"

	"github.com/gin-gonic/gin"
)

// GetUsers godoc
// @Summary Get all users
// @Description Get all active users
// @Tags User
// @Produce  json
// @Security BearerAuth
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

// ChangePassword godoc
// @Summary Change user password
// @Description Change password for the current user
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   body  body  object{old_password=string,new_password=string}  true  "Password change payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /users/change-password [post]
func ChangePassword(c *gin.Context) {
	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Ambil username dari JWT claims (misal sudah ada middleware JWT)
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	username := claims.(*middleware.Claims).Username

	userService := services.UserService{DB: db}
	if err := userService.ChangePassword(username, input.OldPassword, input.NewPassword); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Password changed successfully"})
}
