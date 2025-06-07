// @title Service System REST API
// @version 1.0
// @description REST API for user management and authentication
// @host localhost:8800
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"service_system/controllers"
	_ "service_system/docs"
	"service_system/logger"
	"service_system/router"
	"service_system/utils"

	"github.com/gin-gonic/gin"

	"os"
)

func main() {
	logger.InitLogger()

	utils.Init()
	logger.Log.Info("GIN_MODE:", os.Getenv("GIN_MODE"))
	gin.SetMode(os.Getenv("GIN_MODE"))
	controllers.SetDB(utils.DB)

	r := router.SetupRouter()

	port := os.Getenv("PORT")
	logger.Log.Info("Server running on port:", port)
	if err := r.Run(":" + port); err != nil {
		logger.Log.Fatal("Failed to start server:", err)
	}
}
