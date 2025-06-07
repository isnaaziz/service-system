package controllers

import (
	"fmt"
	"service_system/services"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDB(database *gorm.DB) {
	db = database
}

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if db == nil {
		c.JSON(500, gin.H{"error": "Database not initialized"})
		return
	}
	authService := services.AuthService{DB: db}
	user, err := authService.Register(input.Username, input.Password, input.Email)
	if err != nil {
		fmt.Println("Register error:", err) // log ke terminal
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User registered successfully", "user": user})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if db == nil {
		c.JSON(500, gin.H{"error": "Database not initialized"})
		return
	}
	authService := services.AuthService{DB: db}
	token, err := authService.Login(input.Username, input.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func Logout(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	token := c.GetHeader("Authorization")
	if token == "" {
		// Coba ambil dari body JSON
		if err := c.ShouldBindJSON(&req); err != nil || req.Token == "" {
			c.JSON(400, gin.H{"error": "No token provided"})
			return
		}
		token = req.Token
	}

	// Hapus prefix "Bearer " jika ada
	token = strings.TrimPrefix(token, "Bearer ")

	authService := services.AuthService{DB: db}
	if err := authService.Logout(token); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Logout successful"})
}
