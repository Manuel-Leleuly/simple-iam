package helpers

import (
	"os"
	"time"

	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/golang-jwt/jwt/v5"
)

// Access Token
func CreateAccessToken(user models.User) (tokenString string, err error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	return accessToken.SignedString([]byte(os.Getenv("CLIENT_SECRET")))
}

// Refresh Token
func CreateRefreshToken(accessToken string) (tokenString string, err error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": accessToken,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	return refreshToken.SignedString([]byte(os.Getenv("CLIENT_SECRET")))
}
