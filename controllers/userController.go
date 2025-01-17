package controllers

import (
	"net/http"

	"github.com/Manuel-Leleuly/simple-iam/initializers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	// Get data from request body
	var reqBody models.UserRequest

	if err := c.Bind(&reqBody); err != nil {
		userCreationErrorMessage(c)
		return
	}

	// check if email is already used
	var user models.User
	initializers.DB.First(&user, "email = ?", reqBody.Email)

	if user.Id != "" {
		c.JSON(http.StatusBadRequest, models.ErrorMessage{
			Message: "Email already used",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		userCreationErrorMessage(c)
		return
	}

	// Create the newUser
	newUser := models.User{
		Name:     reqBody.Name,
		Username: reqBody.Username,
		Email:    reqBody.Email,
		Password: string(hash),
	}
	result := initializers.DB.Create(&newUser)

	if result.Error != nil {
		userCreationErrorMessage(c)
		return
	}

	// Send the result
	c.JSON(http.StatusOK, newUser)
}

// helpers
func userCreationErrorMessage(c *gin.Context) {
	c.JSON(http.StatusBadRequest, models.ErrorMessage{
		Message: "Failed to create user",
	})
}
