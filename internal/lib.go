package lib

import (
	"genuinebnt/newsletter/api"

	"github.com/gin-gonic/gin"
)

func Server() *gin.Engine {
	server := gin.Default()

	server.GET("/health_check", api.HealthCheck)
	server.POST("/subscriptions", api.Subscribe)

	return server
}
