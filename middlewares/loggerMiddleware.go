package middlewares

import (
	"time"

	"github.com/Manuel-Leleuly/simple-iam/helpers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(c *gin.Context) {
	logger := helpers.NewLogger()

	startTime := time.Now()

	c.Next()

	latency := time.Since(startTime)

	logger.WithFields(logrus.Fields{
		"METHOD":    c.Request.Method,
		"URI":       c.Request.RequestURI,
		"STATUS":    c.Writer.Status(),
		"LATENCY":   latency,
		"CLIENT_IP": c.ClientIP(),
	}).Info("HTTP REQUEST")

	c.Next()
}
