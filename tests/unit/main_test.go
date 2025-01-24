package unit

import (
	"testing"

	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/Manuel-Leleuly/simple-iam/models"
)

var D *models.DBInstance = helpers.NewDBClient()

func TestMain(m *testing.M) {
	// connect to test db
	if !D.IsDBConnected() {
		if err := helpers.ConnectToTestDB(D); err != nil {
			panic("[Error] failed to connect to test db due to: " + err.Error())
		}
	}

	// delete all leftover test users
	if err := helpers.DeleteAllTestUsers(D); err != nil {
		panic("[Error] failed to delete all test user due to: " + err.Error())
	}

	// create test user
	if err := helpers.CreateTestUser(D); err != nil {
		panic("[Error] failed to create test user due to: " + err.Error())
	}

	m.Run()

	// delete all test users after the test is done
	if err := helpers.DeleteAllTestUsers(D); err != nil {
		panic("[Error] failed to delete all test user due to: " + err.Error())
	}
}
