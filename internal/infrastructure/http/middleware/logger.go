package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luthfiarsyad/mms/internal/infrastructure/logger"
)

// RequestLogger is a Gin middleware that logs HTTP requests using Zerolog.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// process request
		c.Next()

		// after request
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		logger.L.Info().
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Str("ip", clientIP).
			Dur("latency", duration).
			Msg("HTTP request")
	}
}