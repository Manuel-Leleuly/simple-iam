package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Manuel-Leleuly/simple-iam/initializers"
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAccessToken(c *gin.Context) {
	// Get access token from header
	accessToken := c.GetHeader("Authorization")

	// validate the token
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header)
		}
		return []byte(os.Getenv("CLIENT_SECRET")), nil
	})
	if err != nil {
		authErrorMessage(c)
		return
	}

	// continue if token is valid, abort if otherwise
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check token expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			authErrorMessage(c)
			return
		}

		// find the user with the same id as the id stored in token
		var user models.User
		initializers.DB.First(&user, "id = ? AND email = ?", claims["id"], claims["email"])

		if user.Id == "" {
			authErrorMessage(c)
			return
		}

		// continue
		c.Next()
	} else {
		authErrorMessage(c)
		return
	}
}

// helpers
func authErrorMessage(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, models.ErrorMessage{
		Message: "Unauthorized access",
	})
}
