package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func Logger(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// attach logger to context (optional but powerful)
		reqLogger := log.With().
			Str("request_id", requestID).
			Logger()

		c.Set("logger", &reqLogger)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		event := reqLogger.With().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", status).
			Dur("latency", latency).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Logger()

		// log level based on status
		switch {
		case status >= 500:
			event.Error().Msg("request failed")
		case status >= 400:
			event.Warn().Msg("client error")
		default:
			event.Info().Msg("request")
		}
	}
}
