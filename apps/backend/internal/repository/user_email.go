package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/devrapture/stack-trace/internal/models"
	"github.com/devrapture/stack-trace/internal/users"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserEmailRepository interface {
	Create(ctx context.Context, email *models.UserEmail) error
	GetByNormalizedEmail(ctx context.Context, email string) (*models.UserEmail, error)
	GetPrimaryByUserID(ctx context.Context, userID uuid.UUID) (*models.UserEmail, error)
	MarkVerified(ctx context.Context, emailID uuid.UUID, verifiedAt time.Time) error
}

type userEmailRepo struct {
	db *gorm.DB
}

func NewUserEmailRepository(db *gorm.DB) UserEmailRepository {
	return &userEmailRepo{
		db: db,
	}
}

func (r *userEmailRepo) Create(ctx context.Context, email *models.UserEmail) error {
	display, normalized, err := users.NormalizeEmail(email.Email)
	if err != nil {
		return fmt.Errorf("failed to normalize email: %w", err)
	}
	email.Email = display
	email.NormalizedEmail = normalized
	result := r.db.WithContext(ctx).Create(email)
	if result.Error != nil {
		return fmt.Errorf("failed to create user email: %w", result.Error)
	}
	return nil
}

func (r *userEmailRepo) GetByNormalizedEmail(ctx context.Context, email string) (*models.UserEmail, error) {
	var userEmail models.UserEmail
	_, normalized, err := users.NormalizeEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize email: %w", err)
	}

	result := r.db.WithContext(ctx).First(&userEmail, "normalized_email = ?", normalized)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user email: %w", result.Error)
	}
	return &userEmail, nil
}

func (r *userEmailRepo) GetPrimaryByUserID(ctx context.Context, userID uuid.UUID) (*models.UserEmail, error) {
	var userEmail models.UserEmail
	result := r.db.WithContext(ctx).First(&userEmail, "user_id = ? AND is_primary = true", userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}
	return &userEmail, nil
}

func (r *userEmailRepo) MarkVerified(ctx context.Context, emailID uuid.UUID, verifiedAt time.Time) error {
	result := r.db.WithContext(ctx).Model(&models.UserEmail{}).Where("id = ?", emailID).Updates(map[string]interface{}{
		"verified":    true,
		"verified_at": verifiedAt,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to verify email: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with id %s", emailID)
	}
	return nil
}
