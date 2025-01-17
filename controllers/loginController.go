package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Manuel-Leleuly/simple-iam/initializers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	initializers.DB.First(&user, "email = ?", reqBody.Email)

	if user.Id == "" {
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
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the tokens
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("CLIENT_SECRET")))
	if err != nil {
		failedLoginJson(c)
		return
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("CLIENT_SECRET")))
	if err != nil {
		failedLoginJson(c)
		return
	}

	// set as cookies (TBD)
	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("access_token", accessTokenString, 3600, "", "", false, true)
	// c.SetCookie("refresh_token", refreshTokenString, 3600, "", "", false, true)

	// send the result
	c.JSON(http.StatusOK, models.TokenResponse{
		Status:       "success",
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	})
}

// helper
func failedLoginJson(c *gin.Context) {
	c.JSON(http.StatusBadRequest, models.ErrorMessage{
		Message: "Invalid email and/or password",
	})
}
