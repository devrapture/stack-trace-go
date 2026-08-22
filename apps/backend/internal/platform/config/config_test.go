package config_test

import (
	"os"
	"testing"

	"github.com/devrapture/stack-trace/internal/platform/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

//nolint:gosec // test-only URL with fake credentials
const testDatabaseURL = "postgres://user:example@localhost:5432/dbname"

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantErr    bool
		wantPort   string
		wantAppEnv string
		wantDBURL  string
	}{
		{
			name: "defaults are valid",
			env: map[string]string{
				"DATABASE_URL": testDatabaseURL,
			},
			wantAppEnv: "development",
			wantPort:   "8080",
			wantDBURL:  testDatabaseURL,
		},
		{
			name: "invalid app env",
			env: map[string]string{
				"APP_ENV":      "invalid",
				"DATABASE_URL": testDatabaseURL,
			},
			wantErr: true,
		},
		{
			name:    "missing database url",
			env:     map[string]string{},
			wantErr: true,
		},
		{
			name: "valid overrides",
			env: map[string]string{
				"APP_ENV":      "production",
				"PORT":         "8081",
				"DATABASE_URL": testDatabaseURL,
			},
			wantAppEnv: "production",
			wantPort:   "8081",
			wantDBURL:  testDatabaseURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", "")
			t.Setenv("APP_ENV", "")
			t.Setenv("DATABASE_URL", "")
			for key, val := range tt.env {
				t.Setenv(key, val)
			}

			cfg, err := config.Load()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NotNil(t, cfg)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantAppEnv, cfg.AppEnv)
			assert.Equal(t, tt.wantPort, cfg.Port)
			assert.Equal(t, tt.wantDBURL, cfg.DatabaseURL)
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		wantProd bool
		wantErr  bool
	}{
		{
			name:     "return true for production",
			env:      map[string]string{"APP_ENV": "production", "DATABASE_URL": testDatabaseURL},
			wantProd: true,
		},
		{
			name:     "return false for development",
			env:      map[string]string{"APP_ENV": "development", "DATABASE_URL": testDatabaseURL},
			wantProd: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", "")
			t.Setenv("APP_ENV", "")
			t.Setenv("DATABASE_URL", "")

			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg, err := config.Load()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.wantProd, cfg.IsProduction())
		})
	}
}
