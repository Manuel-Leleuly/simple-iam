package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/Manuel-Leleuly/simple-iam/routes"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	router := routes.GetRoutes(D)

	reqBody := models.Login{
		Email:    helpers.TEST_USER.Email,
		Password: helpers.TEST_USER.Password,
	}

	loginJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/iam/v1/login", strings.NewReader(string(loginJson)))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}
