package router

import (
	"os"
	"service_system/controllers"
	"service_system/utils"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.SetTrustedProxies(nil)

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
	user.Use(utils.JWTAuthMiddleware())
	{
		user.GET("", controllers.GetUsers)
		user.DELETE("/:id", controllers.SoftDeleteUser)
		user.POST("change-password", controllers.ChangePassword)
	}

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
