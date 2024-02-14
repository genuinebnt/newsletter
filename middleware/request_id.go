package middleware

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := uuid.NewRandom()
		// Set request ID header in response for client tracking
		c.Header("X-Request-ID", requestID.String())

		// Initialize Zerolog logger
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Str("request_id", requestID.String()).Logger()

		// Attach logger to context
		c.Set("logger", &logger)
		c.Next()
	}
}
