package middlewares

import (
	"errors"
	"net/http"

	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
)

func CheckAccessToken(d *models.DBInstance, c *gin.Context) (statusCode int, err error) {
	// Get access token from header
	bearerToken := c.GetHeader("Authorization")

	// TODO: find a better way to get access token
	accessToken, err := helpers.GetTokenStringFromHeader(bearerToken)
	if err != nil {
		return authError()
	}

	// validate the token
	if err := helpers.ValidateAccessToken(d, accessToken, false); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: err.Error(),
		})
	} else {
		c.Next()
	}

	// arbitrary return
	return http.StatusOK, nil
}

func CheckRefreshToken(d *models.DBInstance, c *gin.Context) (statusCode int, err error) {
	// Get refresh token form header
	bearerToken := c.GetHeader("Authorization")

	// TODO: find a better way to get bearer token
	refreshToken, err := helpers.GetTokenStringFromHeader(bearerToken)
	if err != nil {
		return authError()
	}

	// validate the token
	if err := helpers.ValidateRefreshToken(d, refreshToken); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorMessage{
			Message: err.Error(),
		})
	} else {
		c.Next()
	}

	return http.StatusOK, nil
}

// helpers
func authError() (statusCode int, err error) {
	return http.StatusUnauthorized, errors.New("unauthorized access")
}
