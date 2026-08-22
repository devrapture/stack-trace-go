package repository

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/devrapture/stack-trace/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestUserEmailRepositoryCreateConcurrentFirstEmails(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is required for repository integration tests")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, sqlDB.Close())
	})

	var indexDefinition string
	err = db.Raw(
		`SELECT indexdef FROM pg_indexes WHERE schemaname = current_schema() AND tablename = 'user_emails' AND indexname = ?`,
		primaryEmailConstraint,
	).Scan(&indexDefinition).Error
	require.NoError(t, err)
	require.Contains(t, strings.ToUpper(indexDefinition), "CREATE UNIQUE INDEX")
	require.Contains(t, strings.ToLower(indexDefinition), "where (is_primary = true)")

	user := &models.User{DisplayName: "Concurrent email test"}
	require.NoError(t, db.Create(user).Error)
	t.Cleanup(func() {
		require.NoError(t, db.Unscoped().Where("user_id = ?", user.ID).Delete(&models.UserEmail{}).Error)
		require.NoError(t, db.Unscoped().Delete(&models.User{}, "id = ?", user.ID).Error)
	})

	const emailCount = 8
	testCtx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)
	var primaryInsertsReady atomic.Int32
	releasePrimaryInserts := make(chan struct{})
	const callbackName = "test:coordinate_primary_email_inserts"
	require.NoError(t, db.Callback().Create().Before("gorm:create").Register(callbackName, func(tx *gorm.DB) {
		email, ok := tx.Statement.Dest.(*models.UserEmail)
		if !ok || !email.IsPrimary {
			return
		}
		if primaryInsertsReady.Add(1) == emailCount {
			close(releasePrimaryInserts)
		}
		select {
		case <-releasePrimaryInserts:
		case <-testCtx.Done():
			tx.AddError(testCtx.Err())
		}
	}))
	t.Cleanup(func() {
		require.NoError(t, db.Callback().Create().Remove(callbackName))
	})

	repo := NewUserEmailRepository(db)
	start := make(chan struct{})
	errs := make(chan error, emailCount)
	var wg sync.WaitGroup
	for i := range emailCount {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			errs <- repo.Create(testCtx, &models.UserEmail{
				UserID: user.ID,
				Email:  fmt.Sprintf("concurrent-%d-%s@example.com", i, uuid.NewString()),
			})
		}(i)
	}

	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}

	var createdCount int64
	require.NoError(t, db.Model(&models.UserEmail{}).Where("user_id = ?", user.ID).Count(&createdCount).Error)
	assert.EqualValues(t, emailCount, createdCount)

	var primaryCount int64
	require.NoError(t, db.Model(&models.UserEmail{}).Where("user_id = ? AND is_primary = true", user.ID).Count(&primaryCount).Error)
	assert.EqualValues(t, 1, primaryCount)
}
