package unit

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/Manuel-Leleuly/simple-iam/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine = routes.GetRoutes(D)

func TestGetAllUser(t *testing.T) {
	tokenData, err := helpers.GetTestToken(D)
	assert.Nil(t, err)

	request := getHTTPRequest(http.MethodGet, "/iam/v1/users", nil, tokenData.AccessToken)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.WithPagination[[]models.User]
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(responseBody.Data))

	assert.Equal(t, helpers.TEST_USER.Id, responseBody.Data[0].Id)
	assert.Equal(t, helpers.TEST_USER.Name.FirstName, responseBody.Data[0].Name.FirstName)
	assert.Equal(t, helpers.TEST_USER.Name.LastName, responseBody.Data[0].Name.LastName)
	assert.Equal(t, helpers.TEST_USER.Username, responseBody.Data[0].Username)
	assert.Equal(t, helpers.TEST_USER.Email, responseBody.Data[0].Email)

	// since it currently only has 1 data, `next` and `prev` should be empty
	assert.Equal(t, 0, len(responseBody.Paging.Next))
	assert.Equal(t, 0, len(responseBody.Paging.Prev))
}

func TestGetAllUserWithQuery(t *testing.T) {
	tokenData, err := helpers.GetTestToken(D)
	assert.Nil(t, err)

	// test first name
	request := getHTTPRequest(http.MethodGet, "/iam/v1/users?firstName="+helpers.TEST_USER.FirstName, nil, tokenData.AccessToken)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.WithPagination[[]models.User]
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(responseBody.Data))
	assert.Equal(t, helpers.TEST_USER.FirstName, responseBody.Data[0].FirstName)

	// test last name
	request = getHTTPRequest(http.MethodGet, "/iam/v1/users?lastName="+helpers.TEST_USER.LastName, nil, tokenData.AccessToken)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response = recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(responseBody.Data))
	assert.Equal(t, helpers.TEST_USER.LastName, responseBody.Data[0].LastName)

	// email
	request = getHTTPRequest(http.MethodGet, "/iam/v1/users?email="+helpers.TEST_USER.Email, nil, tokenData.AccessToken)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response = recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(responseBody.Data))
	assert.Equal(t, helpers.TEST_USER.Email, responseBody.Data[0].Email)
}

func TestGetUserByIdSuccess(t *testing.T) {
	tokenData, err := helpers.GetTestToken(D)
	assert.Nil(t, err)

	request := getHTTPRequest(http.MethodGet, "/iam/v1/users/"+helpers.TEST_USER.Id, nil, tokenData.AccessToken)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.Response[models.User]
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, helpers.TEST_USER.Id, responseBody.Data.Id)
	assert.Equal(t, helpers.TEST_USER.Name.FirstName, responseBody.Data.Name.FirstName)
	assert.Equal(t, helpers.TEST_USER.Name.LastName, responseBody.Data.Name.LastName)
	assert.Equal(t, helpers.TEST_USER.Username, responseBody.Data.Username)
	assert.Equal(t, helpers.TEST_USER.Email, responseBody.Data.Email)
}

func TestGetUserByIdFailed(t *testing.T) {
	tokenData, err := helpers.GetTestToken(D)
	assert.Nil(t, err)

	selectedId := "wrongid123"
	request := getHTTPRequest(http.MethodGet, "/iam/v1/users/"+selectedId, nil, tokenData.AccessToken)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.ErrorMessage
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, "failed to get user "+selectedId, responseBody.Message)
}

func TestUpdateUserSuccess(t *testing.T) {
	tokenData, err := helpers.GetTestToken(D)
	assert.Nil(t, err)

	// sleep to make sure user creation and update doesn't occur at nearly the same time
	// since we check `createdAt` and `updatedAt` at the second not microsecond
	time.Sleep(5 * time.Second)

	// update first name
	reqBody := models.UserUpdateRequest{
		FirstName: "newFirstName",
	}

	updateJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	request := getHTTPRequest(http.MethodPatch, "/iam/v1/users/"+helpers.TEST_USER.Id, strings.NewReader(string(updateJson)), tokenData.AccessToken)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.Response[models.User]
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	// only first name should change
	assert.Equal(t, "newFirstName", responseBody.Data.FirstName)
	assert.Equal(t, helpers.TEST_USER.LastName, responseBody.Data.LastName)
	assert.Equal(t, helpers.TEST_USER.Username, responseBody.Data.Username)

	// `updatedAt` should be later than `createdAt`
	assert.Greater(t, responseBody.Data.UpdatedAt.Unix(), responseBody.Data.CreatedAt.Unix())

	// update last name
	reqBody = models.UserUpdateRequest{
		LastName: "newLastName",
	}

	updateJson, err = json.Marshal(reqBody)
	assert.Nil(t, err)

	request = getHTTPRequest(http.MethodPatch, "/iam/v1/users/"+helpers.TEST_USER.Id, strings.NewReader(string(updateJson)), tokenData.AccessToken)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response = recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	// only last name should change
	assert.Equal(t, "newFirstName", responseBody.Data.FirstName)
	assert.Equal(t, "newLastName", responseBody.Data.LastName)
	assert.Equal(t, helpers.TEST_USER.Username, responseBody.Data.Username)

	// update username
	reqBody = models.UserUpdateRequest{
		Username: "newUsername123",
	}

	updateJson, err = json.Marshal(reqBody)
	assert.Nil(t, err)

	request = getHTTPRequest(http.MethodPatch, "/iam/v1/users/"+helpers.TEST_USER.Id, strings.NewReader(string(updateJson)), tokenData.AccessToken)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response = recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	// only username should change
	assert.Equal(t, "newFirstName", responseBody.Data.FirstName)
	assert.Equal(t, "newLastName", responseBody.Data.LastName)
	assert.Equal(t, "newUsername123", responseBody.Data.Username)
}

func TestUpdateUserFailed(t *testing.T) {
	tokenData, err := helpers.GetTestToken(D)
	assert.Nil(t, err)

	// sleep to make sure user creation and update doesn't occur at nearly the same time
	// since we check `createdAt` and `updatedAt` at the second not microsecond
	time.Sleep(5 * time.Second)

	// first name is too long
	reqBody := models.UserUpdateRequest{
		FirstName: "firstNameIsTooLongItShouldOnlyBeBetweenTwoAndFifteenCharacters",
	}

	updateJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	request := getHTTPRequest(http.MethodPatch, "/iam/v1/users/"+helpers.TEST_USER.Id, strings.NewReader(string(updateJson)), tokenData.AccessToken)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.ValidationErrorMessage
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, "FirstName must be alpha, min length 2, max length 15", responseBody.Message[0])
}

func TestDeleteUserFailed(t *testing.T) {
	tokenData, err := helpers.GetTestToken(D)
	assert.Nil(t, err)

	request := getHTTPRequest(http.MethodDelete, "/iam/v1/users/wrongUserId", nil, tokenData.AccessToken)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.ErrorMessage
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, "user not found for id wrongUserId", responseBody.Message)
}

func TestDeleteUserSuccess(t *testing.T) {
	tokenData, err := helpers.GetTestToken(D)
	assert.Nil(t, err)

	request := getHTTPRequest(http.MethodDelete, "/iam/v1/users/"+helpers.TEST_USER.Id, nil, tokenData.AccessToken)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var responseBody models.Response[string]
	err = json.Unmarshal(body, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, "success", responseBody.Data)
}

// helpers
func getHTTPRequest(method string, path string, body io.Reader, accessToken string) *http.Request {
	request := httptest.NewRequest(method, path, body)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+accessToken)

	return request
}
