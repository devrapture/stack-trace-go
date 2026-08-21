package repository

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/devrapture/stack-trace/internal/errors"
	"github.com/devrapture/stack-trace/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("%s: %w", apperrors.ErrFailedToCreateUser, result.Error)
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, "id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}
	return &user, nil
}
