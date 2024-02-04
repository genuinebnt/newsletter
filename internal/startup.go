package lib

import (
	"genuinebnt/newsletter/api"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func Server(conn *pgx.Conn) *gin.Engine {
	server := gin.Default()

	server.Use(func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	})

	server.GET("/health_check", api.HealthCheck)
	server.POST("/subscriptions", api.Subscribe)

	return server
}
