package controllers

import (
	"net/http"

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
	initializers.DB.First(&user, "email = ?", reqBody.Email)

	if user.Id == "" {
		failedLoginJson(c)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		failedLoginJson(c)
		return
	}

	// Generate token
}

// helper
func failedLoginJson(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Invalid email or password",
	})
}
