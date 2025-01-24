package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/Manuel-Leleuly/simple-iam/routes"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUser(t *testing.T) {
	tokenData, err := helpers.GetTestToken(D)
	assert.Nil(t, err)

	router := routes.GetRoutes(D)

	request := httptest.NewRequest(http.MethodGet, "/iam/v1/users", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+tokenData.AccessToken)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}
