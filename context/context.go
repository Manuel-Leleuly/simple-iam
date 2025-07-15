package context

import (
	"github.com/Manuel-Leleuly/simple-iam/models"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) *models.User {
	contextUser, exists := c.Get("user")
	if !exists {
		return &models.User{}
	}

	user, ok := contextUser.(*models.User)
	if !ok {
		return &models.User{}
	}

	return user
}
