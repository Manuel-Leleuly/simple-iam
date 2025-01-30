package routes

import (
	"github.com/Manuel-Leleuly/simple-iam/controllers"
	"github.com/Manuel-Leleuly/simple-iam/middlewares"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const BasePath string = "/iam/v1"

func GetRoutes(d *models.DBInstance) *gin.Engine {
	router := gin.New()

	// use custom logger but keep the default recovery
	router.Use(middlewares.LoggerMiddleware, gin.Recovery())

	// swagger route
	router.GET("/iam/apidocs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	withoutTokenChecker := router.Group(BasePath, d.CheckDBConnection)
	withoutTokenChecker.POST("/login", d.MakeHTTPHandleFunc(controllers.Login))
	withoutTokenChecker.POST("/users", d.MakeHTTPHandleFunc(controllers.CreateUser))

	withAccessChecker := router.Group(BasePath, d.CheckDBConnection, d.MakeHTTPHandleFunc(middlewares.CheckAccessToken))
	withAccessChecker.GET("/users", d.MakeHTTPHandleFunc(controllers.GetUserList))
	withAccessChecker.GET("/users/:userId", d.MakeHTTPHandleFunc(controllers.GetUserDetail))
	withAccessChecker.PATCH("/users/:userId", d.MakeHTTPHandleFunc(controllers.UpdateUser))
	withAccessChecker.DELETE("/users/:userId", d.MakeHTTPHandleFunc(controllers.DeleteUser))

	withRefreshChecker := router.Group(BasePath, d.CheckDBConnection, d.MakeHTTPHandleFunc(middlewares.CheckRefreshToken))
	withRefreshChecker.GET("/token/refresh", d.MakeHTTPHandleFunc(controllers.RefreshToken))

	return router
}
