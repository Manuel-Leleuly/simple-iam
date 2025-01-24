package controllers

import (
	"errors"
	"net/http"

	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login			godoc
// @Summary 		Login
// @Description 	Login
// @Tags			Auth
// @Router			/iam/v1/login [post]
// @Accept			json
// @Produce			json
// @Param			requestBody	body		models.Login{}		true	"Request Body"
// @Success			200			{object}	models.Response[models.TokenResponse]{}
// @Failure			400			{object}	models.ErrorMessage{}
func Login(d *models.DBInstance, c *gin.Context) (statusCode int, err error) {
	// Get the email and password from request body
	var reqBody models.Login

	if err := c.Bind(&reqBody); err != nil {
		return failedLoginError()
	}

	// Find the user
	var user models.User
	result := d.DB.First(&user, "email = ?", reqBody.Email)

	if result.Error != nil || user.Id == "" {
		return failedLoginError()
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		return failedLoginError()
	}

	// Generate tokens
	accessTokenString, err := helpers.CreateAccessToken(user)
	if err != nil {
		return failedLoginError()
	}

	refreshTokenString, err := helpers.CreateRefreshToken(accessTokenString)
	if err != nil {
		return failedLoginError()
	}

	// send the result
	c.JSON(http.StatusOK, models.Response[models.TokenResponse]{
		Data: models.TokenResponse{
			Status:       "success",
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
		},
	})

	return http.StatusOK, nil
}

// helper
func failedLoginError() (statusCode int, err error) {
	return http.StatusBadRequest, errors.New("invalid email and/or password")
}
