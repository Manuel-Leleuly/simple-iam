package routes

import (
	"net/http"

	"github.com/Manuel-Leleuly/simple-iam/controllers"
	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/Manuel-Leleuly/simple-iam/middlewares"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRoutes(d *models.DBInstance) *gin.Engine {
	router := gin.New()

	// use custom logger but keep the default recovery
	router.Use(middlewares.LoggerMiddleware, gin.Recovery())

	// swagger route
	router.GET("/iam/apidocs/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/iam/apidocs/" {
			c.Redirect(http.StatusFound, "/iam/apidocs/index.html")
		}
		ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL(helpers.GetBaseUrl(c)+"/iam/apidocs/doc.json"))(c)
	})

	v1 := router.Group("/iam/v1", d.CheckDBConnection)
	{
		v1.POST("/login", d.MakeHTTPHandleFunc(controllers.Login))
		v1.POST("/users", d.MakeHTTPHandleFunc(controllers.CreateUser))
	}

	withAccessChecker := v1.Group("/", d.MakeHTTPHandleFunc(middlewares.CheckAccessToken))
	{
		withAccessChecker.GET("/users", d.MakeHTTPHandleFunc(controllers.GetUserList))
		withAccessChecker.GET("/users/:userId", d.MakeHTTPHandleFunc(controllers.GetUserDetail))
		withAccessChecker.PATCH("/users/:userId", d.MakeHTTPHandleFunc(controllers.UpdateUser))
		withAccessChecker.DELETE("/users/:userId", d.MakeHTTPHandleFunc(controllers.DeleteUser))
	}

	withRefreshChecker := v1.Group("/", d.MakeHTTPHandleFunc(middlewares.CheckRefreshToken))
	{
		withRefreshChecker.GET("/token/refresh", d.MakeHTTPHandleFunc(controllers.RefreshToken))
	}

	return router
}
