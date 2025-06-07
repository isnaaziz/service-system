package router

import (
	"service_system/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	auth := r.Group("/")
	{
		auth.POST("signup", controllers.Register)
		auth.POST("login", controllers.Login)
		auth.POST("logout", controllers.Logout)
	}
}
