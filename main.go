package main

import (
	"github.com/Manuel-Leleuly/simple-iam/controllers"
	"github.com/Manuel-Leleuly/simple-iam/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	server := gin.Default()

	server.POST("/login", controllers.Login)

	server.POST("/users", controllers.CreateUser)

	server.Run()
}
