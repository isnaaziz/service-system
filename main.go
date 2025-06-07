// @title Service System REST API
// @version 1.0
// @description REST API for user management and authentication
// @host localhost:8800
// @BasePath /
package main

import (
	"service_system/controllers"
	"service_system/database"
	_ "service_system/docs" // penting untuk swag
	"service_system/router"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	database.Init()
	controllers.SetDB(database.DB)

	r := router.SetupRouter()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8800")
}
