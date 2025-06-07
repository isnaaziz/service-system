package router

import (
	"service_system/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

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
