package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/devrapture/stack-trace/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthIdentityRepository interface {
	Create(
		ctx context.Context,
		identity *models.AuthIdentity,
	) error

	GetByUserIDAndKind(
		ctx context.Context,
		userID uuid.UUID,
		kind models.AuthIdentityKind,
	) (*models.AuthIdentity, error)
}

type authIdentityRepo struct {
	db *gorm.DB
}

func NewAuthIdentityRepository(db *gorm.DB) AuthIdentityRepository {
	return &authIdentityRepo{
		db: db,
	}
}

func (r *authIdentityRepo) Create(
	ctx context.Context,
	identity *models.AuthIdentity,
) error {
	err := r.db.WithContext(ctx).Create(identity).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("auth identity already exists: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to create auth identity: %w", err)
	}
	return nil
}

func (r *authIdentityRepo) GetByUserIDAndKind(
	ctx context.Context,
	userID uuid.UUID,
	kind models.AuthIdentityKind,
) (*models.AuthIdentity, error) {
	var userIdentity models.AuthIdentity
	err := r.db.WithContext(ctx).First(&userIdentity, "user_id = ? AND kind = ?", userID, kind).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get auth identity: %w", err)
	}
	return &userIdentity, nil
}
