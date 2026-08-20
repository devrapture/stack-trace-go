package httpserver

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReadinessChecker func(ctx context.Context) error

func HealthzHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func ReadyzHandler(ready ReadinessChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := ready(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	}
}
