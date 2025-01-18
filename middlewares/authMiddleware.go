package middlewares

import (
	"net/http"

	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
)

func CheckAccessToken(c *gin.Context) {
	// Get access token from header
	bearerToken := c.GetHeader("Authorization")

	// TODO: find a better way to get access token
	accessToken, err := helpers.GetTokenStringFromHeader(bearerToken)
	if err != nil {
		authErrorMessage(c)
		return
	}

	// validate the token
	if err := helpers.ValidateAccessToken(accessToken, false); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: err.Error(),
		})
	} else {
		c.Next()
	}
}

func CheckRefreshToken(c *gin.Context) {
	// Get refresh token form header
	bearerToken := c.GetHeader("Authorization")

	// TODO: find a better way to get bearer token
	refreshToken, err := helpers.GetTokenStringFromHeader(bearerToken)
	if err != nil {
		authErrorMessage(c)
		return
	}

	// validate the token
	if err := helpers.ValidateRefreshToken(refreshToken); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: err.Error(),
		})
	} else {
		c.Next()
	}
}

// helpers
func authErrorMessage(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
		Message: "Unauthorized access",
	})
}
