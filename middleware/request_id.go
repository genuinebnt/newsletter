package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := uuid.NewRandom()
		// Attach request ID to context
		c.Set("request_id", requestID)
		// Set request ID header in response for client tracking
		c.Header("X-Request-ID", requestID.String())

		// Initialize Zerolog logger
		logger := zerolog.New(os.Stdout).With().Timestamp().Str("request_id", requestID.String()).Logger()

		// Attach logger to context
		c.Set("logger", &logger)
		c.Next()

	}
}
