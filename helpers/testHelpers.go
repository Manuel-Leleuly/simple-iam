package helpers

import (
	"os"
	"time"

	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var TEST_USER models.User = models.User{
	Id: "d827177441a944cfbed1519635563eaa",
	Name: models.Name{
		FirstName: "Test",
		LastName:  "User",
	},
	Username: "testUser",
	Password: "testing123",
	Email:    "testuser@example.com",
	TimeRecord: models.TimeRecord{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

func ConnectToTestDB(d *models.DBInstance) error {
	err := godotenv.Load("../../.env")
	if err != nil {
		return err
	}

	err = d.ConnectToDB(os.Getenv("DB_TEST_NAME"))
	if err != nil {
		return err
	}

	err = d.SyncDatabase()
	if err != nil {
		return err
	}

	return nil
}

func GetTestToken(d *models.DBInstance) (tokenData *models.TokenResponse, err error) {
	loginData := models.Login{
		Email:    TEST_USER.Email,
		Password: TEST_USER.Password,
	}

	// get user
	var user models.User
	result := d.DB.First(&user, "email = ?", loginData.Email)
	if result.Error != nil || user.Id == "" {
		return nil, result.Error
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		return nil, err
	}

	// generate token
	accessTokenString, err := CreateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := CreateRefreshToken(accessTokenString)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		Status:       "success",
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func CreateTestUser(d *models.DBInstance) error {
	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(TEST_USER.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// create the new user
	newUser := models.User{
		Id:       TEST_USER.Id,
		Name:     TEST_USER.Name,
		Username: TEST_USER.Username,
		Password: string(hash),
		Email:    TEST_USER.Email,
	}
	result := d.DB.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteAllTestUsers(d *models.DBInstance) error {
	var users []models.User
	result := d.DB.Raw("TRUNCATE USERS").Scan(&users)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
