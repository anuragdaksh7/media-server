package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()

		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		requestID := c.GetString(RequestIDKey)

		c.Next()

		latency := time.Since(start)

		statusCode := c.Writer.Status()

		userID := c.GetString("user_id")

		slog.Info(
			"http request",
			slog.String("request_id", requestID),
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.Duration("latency", latency),
			slog.String("ip", ip),
			slog.String("user_agent", userAgent),
			slog.String("user_id", userID),
		)
	}
}
