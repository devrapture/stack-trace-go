package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devrapture/stack-trace/internal/platform/config"
	"github.com/devrapture/stack-trace/internal/platform/database"
	"github.com/devrapture/stack-trace/internal/platform/httpserver"
	"github.com/devrapture/stack-trace/internal/platform/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	logger, err := logger.NewLogger(!cfg.IsProduction())
	if err != nil {
		return fmt.Errorf("failed to create logger: %v", err)
	}
	defer flushLog(logger)
	logger.Info("starting stack trace server", zap.String("APP_ENV", cfg.AppEnv), zap.String("PORT", cfg.Port))

	db, err := database.ConnectDb(cfg)
	if err != nil {
		logger.Error("Failed to initialize database", zap.Error(err))
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to access database connection pool", zap.Error(err))
		return err
	}
	defer func() {
		_ = sqlDB.Close()
	}()

	addr := fmt.Sprintf(":%s", cfg.Port)
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := httpserver.Setup(cfg, db, logger)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r.Handler(),
		ReadHeaderTimeout: 5 * time.Second, // mitigates slowloris-style attacks
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		err := srv.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		serverErrors <- err
	}()

	select {
	case err := <-serverErrors:
		return err
	case <-signalCtx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		_ = srv.Close()
		return fmt.Errorf("shutdown HTTP server:%w", err)
	}

	if err := <-serverErrors; err != nil {
		return fmt.Errorf("HTTP server stopped: %w", err)
	}

	logger.Info("server stopped gracefully")

	return nil
}

func flushLog(logger *zap.Logger) {
	_ = logger.Sync()
}
