package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devrapture/stack-trace/internal/platform/config"
	"github.com/devrapture/stack-trace/internal/platform/httpserver"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//nolint:gosec // test-only URL with fake credentials
const testDatabaseURL = "postgres://user:example@localhost:5432/testdb"

func setupRouter(db *gorm.DB) *gin.Engine {
	cfg := &config.Config{
		AppEnv:      "development",
		Port:        "8080",
		DatabaseURL: testDatabaseURL,
	}

	logger := zap.NewNop()

	gin.SetMode(gin.TestMode)

	return httpserver.Setup(cfg, db, logger)
}

func TestHealth(t *testing.T) {
	router := setupRouter(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestRequestIDHeaderIsSet(t *testing.T) {
	router := setupRouter(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Header().Get(httpserver.RequestIDHeader) == "" {
		t.Fatalf(
			"expected %s header to be set",
			httpserver.RequestIDHeader,
		)
	}
}
