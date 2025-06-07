package router

import (
	"os"
	"service_system/controllers"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Ambil origins dari .env
	origins := os.Getenv("CORS_ORIGINS")
	allowedOrigins := strings.Split(origins, ",")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Auth group
	auth := r.Group("/")
	{
		auth.POST("signup", controllers.Register)
		auth.POST("login", controllers.Login)
		auth.POST("logout", controllers.Logout)
	}

	// User group
	user := r.Group("/users")
	{
		user.GET("", controllers.GetUsers)
		user.DELETE("/:id", controllers.SoftDeleteUser)
	}

	return r
}
