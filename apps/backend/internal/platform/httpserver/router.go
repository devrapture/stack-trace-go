package httpserver

import (
	"context"

	"github.com/devrapture/stack-trace/internal/platform/config"
	"github.com/devrapture/stack-trace/internal/platform/database"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB, logger *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(
		gin.Recovery(),
		RequestID(),
		ZapLogger(logger),
	)

	r.
		GET("/healthz", HealthzHandler()).
		GET("/readyz", ReadyzHandler(func(ctx context.Context) error {
			return database.Ping(ctx, db)
		}))
	// v1 := r.Group("/api/v1")
	// {
	// 	v1.
	// 		GET("/healthz", HealthzHandler()).
	// 		GET("/readyz", ReadyzHandler(func(ctx context.Context) error {
	// 			return database.Ping(ctx, db)
	// 		}))
	// }

	return r
}
