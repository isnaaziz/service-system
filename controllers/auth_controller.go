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

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, password, and email
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   body  body  object{username=string,password=string,email=string}  true  "User credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /signup [post]
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

// Login godoc
// @Summary Login user
// @Description Login with username and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   body  body  object{username=string,password=string}  true  "User credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /login [post]
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

	c.JSON(200, gin.H{
		"message": "Login successfully",
		"token":   token,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout with JWT token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   Authorization header string false "Bearer token"
// @Param   body  body  object{token=string}  false  "JWT token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /logout [post]
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
