package routes

import (
	"github.com/Manuel-Leleuly/simple-iam/controllers"
	"github.com/Manuel-Leleuly/simple-iam/middlewares"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const BasePath string = "/iam/v1"

func GetRoutes() *gin.Engine {
	router := gin.Default()

	// swagger route
	router.GET("/iam/apidocs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	withoutTokenChecker := router.Group(BasePath)
	withoutTokenChecker.POST("/login", controllers.Login)
	withoutTokenChecker.POST("/users", controllers.CreateUser)

	withAccessChecker := router.Group(BasePath, middlewares.CheckAccessToken)
	withAccessChecker.GET("/users", controllers.GetUserList)
	withAccessChecker.GET("/users/:userId", controllers.GetUserDetail)
	withAccessChecker.PATCH("/users/:userId", controllers.UpdateUser)
	withAccessChecker.DELETE("/users/:userId", controllers.DeleteUser)

	withRefreshChecker := router.Group(BasePath, middlewares.CheckRefreshToken)
	withRefreshChecker.GET("/token/refresh", middlewares.CheckRefreshToken, controllers.RefreshToken)

	return router
}
