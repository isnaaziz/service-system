package router

import (
	"service_system/controllers"
	"service_system/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	user := r.Group("/users")
	user.Use(middleware.JWTAuthMiddleware())
	{
		user.GET("", controllers.GetUsers)
		user.DELETE("/:id", controllers.SoftDeleteUser)
		user.POST("change-password", controllers.ChangePassword)
	}
}
