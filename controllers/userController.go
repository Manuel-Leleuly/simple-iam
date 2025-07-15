package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Manuel-Leleuly/simple-iam/constants"
	"github.com/Manuel-Leleuly/simple-iam/helpers"
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
// @Success			200			{object}	models.Response[models.UserResponse]{}
// @Failure			400			{object}	models.ErrorMessage{}
func CreateUser(d *models.DBInstance, c *gin.Context) (statusCode int, err error) {
	// Get data from request body
	var reqBody models.UserRequest

	if err := c.Bind(&reqBody); err != nil {
		return createUserFailed()
	}

	// validate request body
	validate := validation.GetValidator()
	if err := validate.Struct(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorMessage{
			Message: validation.TranslateValidationErrors(err),
		})

		// arbitrary return so that the function stops
		return http.StatusBadRequest, nil
	}

	// check if email is already used
	var user models.User
	d.DB.First(&user, "email = ?", reqBody.Email)

	if user.Id != "" {
		return http.StatusBadRequest, errors.New("email already used")
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return createUserFailed()
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
	result := d.DB.Create(&newUser)

	if result.Error != nil {
		return createUserFailed()
	}

	// Send the result
	c.JSON(http.StatusOK, models.Response[models.UserResponse]{
		Data: helpers.ConvertUserToUserResponse(newUser),
	})

	return http.StatusOK, nil
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
// @Success			200			{object}	models.WithPagination[[]models.UserResponse]{}
// @Failure			400			{object}	models.ErrorMessage{}
func GetUserList(d *models.DBInstance, c *gin.Context) (statusCode int, err error) {
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

	/*
		Add an extra data to check whether there's a next page or not
		this is to avoid calling DB for the 2nd time

		This extra data will NOT be included in the response
	*/
	dbQuery := d.DB.Limit(selectedLimit + 1)

	if len(firstName) > 0 {
		dbQuery = dbQuery.Where("first_name like ?", "%"+firstName+"%")
	}
	if len(lastName) > 0 {
		dbQuery = dbQuery.Where("last_name like ?", "%"+lastName+"%")
	}
	if len(email) > 0 {
		dbQuery = dbQuery.Where("email = ?", email)
	}

	result := dbQuery.Offset(selectedOffset).Find(&users)

	if result.Error != nil {
		return getUserListFailed()
	}

	hasNext := len(users) > selectedLimit+1

	// get paging
	paging, err := helpers.GetPagination(helpers.GetFullUrl(c), hasNext)
	if err != nil {
		return getUserListFailed()
	}

	// return the result
	var userResponses []models.UserResponse
	for index, user := range users {
		if index >= selectedLimit {
			break
		}
		userResponses = append(userResponses, helpers.ConvertUserToUserResponse(user))
	}

	c.JSON(http.StatusOK, models.WithPagination[[]models.UserResponse]{
		Data:   userResponses,
		Paging: *paging,
	})

	return http.StatusOK, nil
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
// @Success			200			{object}	models.Response[models.UserResponse]{}
// @Failure			400			{object}	models.ErrorMessage{}
func GetUserDetail(d *models.DBInstance, c *gin.Context) (statusCode int, err error) {
	// get id param
	idParam := c.Param("userId")

	// get the user
	var user models.User

	result := d.DB.Where("id = ?", idParam).First(&user)
	if result.Error != nil || user.Id == "" {
		return getUserDetailFailed(idParam)
	}

	// return the result
	c.JSON(http.StatusOK, models.Response[models.UserResponse]{
		Data: helpers.ConvertUserToUserResponse(user),
	})

	return http.StatusOK, nil
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
// @Success			200			{object}	models.Response[models.UserResponse]{}
// @Failure			400			{object}	models.ErrorMessage{}
// @Failure			404			{object}	models.ErrorMessage{}
func UpdateUser(d *models.DBInstance, c *gin.Context) (statusCode int, err error) {
	// get request body and param from url
	var reqBody models.UserUpdateRequest
	idParam := c.Param("userId")

	if err := c.Bind(&reqBody); err != nil {
		return updateUserFailed(idParam)
	}

	// validate request body
	validate := validation.GetValidator()
	if err := validate.Struct(reqBody); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorMessage{
			Message: validation.TranslateValidationErrors(err),
		})

		// arbitrary return so that the function stops
		return http.StatusBadRequest, nil
	}

	// get user
	var user models.User

	result := d.DB.Where("id = ?", idParam).First(&user)
	if result.Error != nil || user.Id == "" {
		return http.StatusNotFound, errors.New("user not found for id " + idParam)
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

	result = d.DB.Save(&user)
	if result.Error != nil {
		return updateUserFailed(idParam)
	}

	// return the updated user
	c.JSON(http.StatusOK, models.Response[models.User]{
		Data: user,
	})

	return http.StatusOK, nil
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
func DeleteUser(d *models.DBInstance, c *gin.Context) (statusCode int, err error) {
	// get id param
	idParam := c.Param("userId")

	// get user
	var user models.User

	result := d.DB.Where("id = ?", idParam).First(&user)
	if result.Error != nil || user.Id == "" {
		return http.StatusNotFound, errors.New("user not found for id " + idParam)
	}

	// remove user
	result = d.DB.Delete(&user, "id = ?", user.Id)
	if result.Error != nil {
		return deleteUserFailed(idParam)
	}

	// return response
	c.JSON(http.StatusOK, models.Response[string]{
		Data: "success",
	})

	return http.StatusOK, nil
}

// helpers
func createUserFailed() (statusCode int, err error) {
	return http.StatusBadRequest, errors.New("failed to create user")
}

func getUserListFailed() (statusCode int, err error) {
	return http.StatusBadRequest, errors.New("failed to get users")
}

func getUserDetailFailed(id string) (statusCode int, err error) {
	return http.StatusNotFound, errors.New("failed to get user " + id)
}

func updateUserFailed(id string) (statusCode int, err error) {
	return http.StatusBadRequest, errors.New("failed to update user " + id)
}

func deleteUserFailed(id string) (statusCode int, err error) {
	return http.StatusBadRequest, errors.New("failed to delete user " + id)
}
