package httpserver

import (
	"net/http"

	"github.com/devrapture/stack-trace/internal/platform/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Setup(cfg *config.Config, logger *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(
		gin.Recovery(),
		RequestID(),
		ZapLogger(logger),
	)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "OK"})
		})
	}

	return r
}
