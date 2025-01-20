package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Manuel-Leleuly/simple-iam/constants"
	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/Manuel-Leleuly/simple-iam/initializers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/Manuel-Leleuly/simple-iam/validation"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser		godoc
// @Summary 		Create user
// @Description 	Create a user
// @Tags			User
// @Router			/iam/v1/users [post]
// @Accept			json
// @Produce			json
// @Param			requestBody	body		models.UserRequest{}		true	"Request Body"
// @Success			200			{object}	models.Response[models.User]{}
// @Failure			400			{object}	models.ErrorMessage{}
func CreateUser(c *gin.Context) {
	// Get data from request body
	var reqBody models.UserRequest

	if err := c.Bind(&reqBody); err != nil {
		userCreationErrorMessage(c)
		return
	}

	// validate request body
	validate := validation.GetValidator()
	if err := validate.Struct(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorMessage{
			Message: validation.TranslateValidationErrors(err),
		})
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
		Name: models.Name{
			FirstName: reqBody.FirstName,
			LastName:  reqBody.LastName,
		},
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
	c.JSON(http.StatusOK, gin.H{
		"data": newUser,
	})
}

// GetUserList		godoc
// @Summary 		Get user list
// @Description 	Get all users from database
// @Security 		ApiKeyAuth
// @Tags			User
// @Router			/iam/v1/users [get]
// @Accept			json
// @Produce			json
// @Param			firstName	query 		string		false	"search by first name"
// @Param			lastName	query		string		false	"search by last name"
// @Param			email		query		string		false	"search by email"
// @Param			offset		query		string		false	"default to 0"
// @Param			limit		query		string		false	"default to 10"
// @Success			200			{object}	models.WithPagination[[]models.User]{}
// @Failure			400			{object}	models.ErrorMessage{}
func GetUserList(c *gin.Context) {
	// get all query params
	firstName := c.Query("firstName")
	lastName := c.Query("lastName")
	email := c.Query("email")
	offset := c.Query("offset")
	limit := c.Query("limit")

	selectedOffset, err := strconv.Atoi(offset)
	if err != nil {
		selectedOffset = constants.DEFAULT_OFFSET
	}

	selectedLimit, err := strconv.Atoi(limit)
	if err != nil {
		selectedLimit = constants.DEFAULT_LIMIT
	}

	// get users
	var users []models.User

	dbQuery := initializers.DB.Offset(selectedOffset).Limit(selectedLimit)

	if len(firstName) > 0 {
		dbQuery = dbQuery.Where("first_name like ?", "%"+firstName+"%")
	}
	if len(lastName) > 0 {
		dbQuery = dbQuery.Where("last_name like ?", "%"+lastName+"%")
	}
	if len(email) > 0 {
		dbQuery = dbQuery.Where("email = ?", email)
	}

	result := dbQuery.Find(&users)

	if result.Error != nil {
		getUserListErrorMessage(c)
		return
	}

	// get paging
	paging, err := helpers.GetPagination(helpers.GetFullUrl(c))
	if err != nil {
		fmt.Println(err)
		getUserListErrorMessage(c)
		return
	}

	// return the result
	c.JSON(http.StatusOK, models.WithPagination[[]models.User]{
		Data:   users,
		Paging: *paging,
	})
}

// GetUserDetail	godoc
// @Summary 		Get user detail
// @Description 	Get the details of a user
// @Security 		ApiKeyAuth
// @Tags			User
// @Router			/iam/v1/users/{userId} [get]
// @Accept			json
// @Produce			json
// @Param			userId		path 		string		true	"User ID"
// @Success			200			{object}	models.Response[models.User]{}
// @Failure			400			{object}	models.ErrorMessage{}
func GetUserDetail(c *gin.Context) {
	// get id param
	idParam := c.Param("userId")

	// get the user
	var user models.User

	result := initializers.DB.Where("id = ?", idParam).First(&user)
	if result.Error != nil || user.Id == "" {
		getUserDetailErrorMessage(c, idParam)
		return
	}

	// return the result
	c.JSON(http.StatusOK, models.Response[models.User]{
		Data: user,
	})
}

// UpdateUser		godoc
// @Summary 		Update user
// @Description 	Update the user
// @Security 		ApiKeyAuth
// @Tags			User
// @Router			/iam/v1/users/{userId} [patch]
// @Accept			json
// @Produce			json
// @Param			userId		path 		string							true	"User ID"
// @Param			requestBody	body		models.UserUpdateRequest{}		true	"Request Body"
// @Success			200			{object}	models.Response[models.User]{}
// @Failure			400			{object}	models.ErrorMessage{}
// @Failure			404			{object}	models.ErrorMessage{}
func UpdateUser(c *gin.Context) {
	// get request body and param from url
	var reqBody models.UserUpdateRequest
	idParam := c.Param("userId")

	if err := c.Bind(&reqBody); err != nil {
		fmt.Println("result error 1: ", err)
		updateUserErrorMessage(c, idParam)
		return
	}

	// validate request body
	validate := validation.GetValidator()
	if err := validate.Struct(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorMessage{
			Message: validation.TranslateValidationErrors(err),
		})
		return
	}

	// get user
	var user models.User

	result := initializers.DB.Where("id = ?", idParam).First(&user)
	if result.Error != nil || user.Id == "" {
		c.JSON(http.StatusNotFound, models.ErrorMessage{
			Message: "User not found for id " + idParam,
		})
		return
	}

	// update the user
	if reqBody.FirstName != "" {
		user.FirstName = reqBody.FirstName
	}
	if reqBody.LastName != "" {
		user.LastName = reqBody.LastName
	}
	if reqBody.Username != "" {
		user.Username = reqBody.Username
	}

	result = initializers.DB.Save(&user)
	if result.Error != nil {
		updateUserErrorMessage(c, idParam)
		return
	}

	// return the updated user
	c.JSON(http.StatusOK, models.Response[models.User]{
		Data: user,
	})
}

// DeleteUser		godoc
// @Summary 		Delete user
// @Description 	Delete a user
// @Security 		ApiKeyAuth
// @Tags			User
// @Router			/iam/v1/users/{userId} [delete]
// @Accept			json
// @Produce			json
// @Param			userId		path 		string		true	"User ID"
// @Success			200			{object}	models.Response[string]{}
// @Failure			400			{object}	models.ErrorMessage{}
// @Failure			404			{object}	models.ErrorMessage{}
func DeleteUser(c *gin.Context) {
	// get id param
	idParam := c.Param("userId")

	// get user
	var user models.User

	result := initializers.DB.Where("id = ?", idParam).First(&user)
	if result.Error != nil || user.Id == "" {
		c.JSON(http.StatusNotFound, models.ErrorMessage{
			Message: "User not found for id " + idParam,
		})
		return
	}

	// remove user
	result = initializers.DB.Delete(&user, "id = ?", user.Id)
	if result.Error != nil {
		deleteUserErrorMessage(c, idParam)
		return
	}

	// return response
	c.JSON(http.StatusOK, models.Response[string]{
		Data: "success",
	})
}

// helpers
func userCreationErrorMessage(c *gin.Context) {
	c.JSON(http.StatusBadRequest, models.ErrorMessage{
		Message: "Failed to create user",
	})
}

func getUserListErrorMessage(c *gin.Context) {
	c.JSON(http.StatusBadRequest, models.ErrorMessage{
		Message: "Failed to get users",
	})
}

func getUserDetailErrorMessage(c *gin.Context, id string) {
	c.JSON(http.StatusNotFound, models.ErrorMessage{
		Message: "Failed to get user " + id,
	})
}

func updateUserErrorMessage(c *gin.Context, id string) {
	c.JSON(http.StatusBadRequest, models.ErrorMessage{
		Message: "Failed to update user " + id,
	})
}

func deleteUserErrorMessage(c *gin.Context, id string) {
	c.JSON(http.StatusBadRequest, models.ErrorMessage{
		Message: "Failed to delete user " + id,
	})
}
