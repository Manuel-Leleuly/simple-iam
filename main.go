package main

import (
	"github.com/Manuel-Leleuly/simple-iam/controllers"
	"github.com/Manuel-Leleuly/simple-iam/initializers"
	"github.com/Manuel-Leleuly/simple-iam/middlewares"
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
	server.GET("/users", middlewares.CheckAccessToken, controllers.GetUserList)
	server.GET("/users/:userId", middlewares.CheckAccessToken, controllers.GetUserDetail)
	server.PATCH("/users/:userId", middlewares.CheckAccessToken, controllers.UpdateUser)
	server.DELETE("/users/:userId", middlewares.CheckAccessToken, controllers.DeleteUser)

	// why port 3005? That's a mystery even I can't answer
	err := server.Run(":3005")
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
