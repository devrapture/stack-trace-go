package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/devrapture/stack-trace/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordCredentialRepository interface {
	Create(ctx context.Context, credentials *models.PasswordCredential) error
	GetByAuthIdentityID(ctx context.Context, authentityID uuid.UUID) (*models.PasswordCredential, error)
	UpdateHash(ctx context.Context, authentityID uuid.UUID, passwordHash string) error
}

type passwordCredentialRepository struct {
	db *gorm.DB
}

func NewPasswordCredentialRepository(db *gorm.DB) PasswordCredentialRepository {
	return &passwordCredentialRepository{
		db: db,
	}
}

func (r *passwordCredentialRepository) Create(ctx context.Context, credentials *models.PasswordCredential) error {
	result := r.db.WithContext(ctx).Create(credentials)
	if result.Error != nil {
		return fmt.Errorf("failed to create password credential: %w", result.Error)
	}
	return nil
}

func (r *passwordCredentialRepository) GetByAuthIdentityID(ctx context.Context, authentityID uuid.UUID) (*models.PasswordCredential, error) {
	var credential models.PasswordCredential
	err := r.db.WithContext(ctx).First(&credential, "auth_identity_id = ?", authentityID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf(
			"failed to get password credential: %w",
			err,
		)
	}
	return &credential, nil
}

func (r *passwordCredentialRepository) UpdateHash(ctx context.Context, authentityID uuid.UUID, passwordHash string) error {
	result := r.db.WithContext(ctx).Model(models.PasswordCredential{}).Where("auth_identity_id = ?", authentityID).Update("password_hash", passwordHash)
	if result.Error != nil {
		return fmt.Errorf("failed to update password credential: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("password credential not found")
	}
	return nil
}
