package controllers

import (
	"fmt"
	"net/http"

	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/Manuel-Leleuly/simple-iam/initializers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	// Get the email and password from request body
	var reqBody models.Login

	if err := c.Bind(&reqBody); err != nil {
		failedLoginJson(c)
		return
	}

	// Find the user
	var user models.User
	result := initializers.DB.First(&user, "email = ?", reqBody.Email)

	if result.Error != nil || user.Id == "" {
		failedLoginJson(c)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		fmt.Println(err)
		failedLoginJson(c)
		return
	}

	// Generate tokens
	accessTokenString, err := helpers.CreateAccessToken(user)
	if err != nil {
		failedLoginJson(c)
		return
	}

	refreshTokenString, err := helpers.CreateRefreshToken(accessTokenString)
	if err != nil {
		failedLoginJson(c)
		return
	}

	// send the result
	c.JSON(http.StatusOK, gin.H{
		"data": models.TokenResponse{
			Status:       "success",
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
		},
	})
}

// helper
func failedLoginJson(c *gin.Context) {
	c.JSON(http.StatusBadRequest, models.ErrorMessage{
		Message: "Invalid email and/or password",
	})
}
