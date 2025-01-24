package controllers

import (
	"errors"
	"net/http"

	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RefreshToken(d *models.DBInstance, c *gin.Context) (statusCode int, err error) {
	// get refresh token from header
	bearerToken := c.GetHeader("Authorization")

	refreshTokenString, err := helpers.GetTokenStringFromHeader(bearerToken)
	if err != nil {
		return refreshTokenFailed()
	}

	refreshToken, err := helpers.GetToken(refreshTokenString)
	if err != nil || !refreshToken.Valid {
		return refreshTokenFailed()
	}

	// get access token from refresh token
	refreshClaims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return refreshTokenFailed()
	}

	accessTokenString, ok := (refreshClaims["token"]).(string)
	if !ok {
		return refreshTokenFailed()
	}

	accessToken, err := helpers.GetToken(accessTokenString)
	if err != nil {
		return refreshTokenFailed()
	}

	// get id and email from access token
	accessClaims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok {
		return refreshTokenFailed()
	}

	userId, ok := (accessClaims["id"]).(string)
	if !ok {
		return refreshTokenFailed()
	}

	userEmail, ok := (accessClaims["email"]).(string)
	if !ok {
		return refreshTokenFailed()
	}

	// generate new access and refresh tokens
	newAccessToken, err := helpers.CreateAccessToken(models.User{
		Id:    userId,
		Email: userEmail,
	})
	if err != nil {
		return refreshTokenFailed()
	}

	newRefreshToken, err := helpers.CreateRefreshToken(newAccessToken)
	if err != nil {
		return refreshTokenFailed()
	}

	// return the result
	c.JSON(http.StatusOK, gin.H{
		"data": models.TokenResponse{
			Status:       "success",
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
		},
	})

	return http.StatusOK, nil
}

// helpers
func refreshTokenFailed() (statusCode int, err error) {
	return http.StatusBadRequest, errors.New("invalid refresh token")
}
