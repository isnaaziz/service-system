package main

import (
	"service_system/controllers"
	"service_system/database"
	"service_system/router"
)

func main() {
	database.Init()
	controllers.SetDB(database.DB)

	r := router.SetupRouter()
	r.Run(":8800")
}
