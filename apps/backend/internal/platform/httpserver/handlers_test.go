package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestHealthzHandler(t *testing.T) {
	r := gin.New()
	r.GET("/healthz", HealthzHandler())

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"ok"}`, rec.Body.String())
}

func TestReadyzHandler_Ready(t *testing.T) {
	ready := func(ctx context.Context) error {
		return nil
	}
	r := gin.New()
	r.GET("/readyz", ReadyzHandler(ready))

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ready"}`, rec.Body.String())
}

func TestReadyzHandler_NotReady(t *testing.T) {
	ready := func(ctx context.Context) error {
		return errors.New("database unavailable")
	}
	r := gin.New()
	r.GET("/readyz", ReadyzHandler(ready))

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.JSONEq(t, `{"status":"not ready"}`, rec.Body.String())
}
