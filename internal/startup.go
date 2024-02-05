package lib

import (
	"genuinebnt/newsletter/api"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Server(dbpool *pgxpool.Pool) *gin.Engine {
	server := gin.Default()

	server.Use(func(c *gin.Context) {
		c.Set("db", dbpool)
		c.Next()
	})

	server.GET("/health_check", api.HealthCheck)
	server.POST("/subscriptions", api.Subscribe)

	return server
}
