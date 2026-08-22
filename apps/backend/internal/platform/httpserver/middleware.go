package httpserver

import (
	"time"

	"github.com/devrapture/stack-trace/internal/platform/random"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	RequestIDHeader     = "X-Request-ID"
	requestIDContextKey = "request_id"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			generatedID, err := random.Hex(8)
			if err != nil {
				generatedID = "unknown"
			}
			requestID = generatedID
		}
		c.Set(requestIDContextKey, requestID)
		c.Writer.Header().Set(RequestIDHeader, requestID)
		c.Next()
	}
}

func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process the request
		c.Next()

		requestID, exists := c.Get(requestIDContextKey)
		if !exists {
			requestID = "unknown"
		}

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logFn := logger.Info
		if statusCode >= 500 {
			logFn = logger.Error
		} else if statusCode >= 400 {
			logFn = logger.Warn
		}

		logFn("HTTP Request",
			zap.String("request_id", requestID.(string)),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("duration", duration),
			zap.String("ip", c.ClientIP()),
		)
	}
}
