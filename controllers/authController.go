package controllers

import (
	"net/http"

	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RefreshToken(c *gin.Context) {
	// get refresh token from header
	bearerToken := c.GetHeader("Authorization")

	refreshTokenString, err := helpers.GetTokenStringFromHeader(bearerToken)
	if err != nil {
		refreshTokenErrorMessage(c)
		return
	}

	refreshToken, err := helpers.GetToken(refreshTokenString)
	if err != nil || !refreshToken.Valid {
		refreshTokenErrorMessage(c)
		return
	}

	// get access token from refresh token
	refreshClaims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		refreshTokenErrorMessage(c)
		return
	}

	accessTokenString, ok := (refreshClaims["token"]).(string)
	if !ok {
		refreshTokenErrorMessage(c)
		return
	}

	accessToken, err := helpers.GetToken(accessTokenString)
	if err != nil {
		refreshTokenErrorMessage(c)
		return
	}

	// get id and email from access token
	accessClaims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok {
		refreshTokenErrorMessage(c)
		return
	}

	userId, ok := (accessClaims["id"]).(string)
	if !ok {
		refreshTokenErrorMessage(c)
		return
	}

	userEmail, ok := (accessClaims["email"]).(string)
	if !ok {
		refreshTokenErrorMessage(c)
		return
	}

	// generate new access and refresh tokens
	newAccessToken, err := helpers.CreateAccessToken(models.User{
		Id:    userId,
		Email: userEmail,
	})
	if err != nil {
		refreshTokenErrorMessage(c)
		return
	}

	newRefreshToken, err := helpers.CreateRefreshToken(newAccessToken)
	if err != nil {
		refreshTokenErrorMessage(c)
		return
	}

	// return the result
	c.JSON(http.StatusOK, gin.H{
		"data": models.TokenResponse{
			Status:       "success",
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
		},
	})
}

// helpers
func refreshTokenErrorMessage(c *gin.Context) {
	c.JSON(http.StatusBadRequest, models.ErrorMessage{
		Message: "invalid refresh token",
	})
}
