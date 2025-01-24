package helpers

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

func ValidateAccessToken(d *models.DBInstance, accessToken string, isFromRefreshToken bool) error {
	// get the token
	token, err := GetToken(accessToken)
	if err != nil {
		return err
	}

	// check if token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check token expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) && !isFromRefreshToken {
			return errors.New("token is expired")
		}

		// find the user with the same id as the id stored in token
		var user models.User
		result := d.DB.First(&user, "id = ? AND email = ?", claims["id"], claims["email"])

		if result.Error != nil || user.Id == "" {
			return errors.New("unauthorized access")
		}
	} else {
		return errors.New("unauthorized access")
	}

	return nil
}

// Refresh Token
func CreateRefreshToken(accessToken string) (tokenString string, err error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": accessToken,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	return refreshToken.SignedString([]byte(os.Getenv("CLIENT_SECRET")))
}

func ValidateRefreshToken(d *models.DBInstance, refreshToken string) error {
	// get the token
	token, err := GetToken(refreshToken)
	if err != nil {
		return err
	}

	// check if token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check token expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return errors.New("token is expired")
		}

		// validate the data inside access token
		accessToken := claims["token"]

		accessTokenString, ok := accessToken.(string)
		if !ok {
			return errors.New("unauthorized access")
		}

		if err := ValidateAccessToken(d, accessTokenString, true); err != nil {
			return err
		}
	} else {
		return errors.New("unauthorized access")
	}

	return nil
}

// helpers
func GetToken(tokenString string) (token *jwt.Token, err error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header)
		}
		return []byte(os.Getenv("CLIENT_SECRET")), nil
	})
}

func GetTokenStringFromHeader(bearerToken string) (string, error) {
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		return "", errors.New("invalid bearer token")
	}

	tokenString := strings.ReplaceAll(bearerToken, "Bearer ", "")

	return tokenString, nil
}
