package models

import (
	"time"

	"github.com/google/uuid"
)

type UserEmail struct {
	Base

	UserID          uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:uq_user_emails_primary,where:is_primary = true AND deleted_at IS NULL"`
	Email           string     `json:"email" gorm:"not null;"`
	NormalizedEmail string     `json:"-" gorm:"not null;uniqueIndex"`
	IsPrimary       bool       `json:"is_primary" gorm:"not null"`
	Verified        bool       `json:"verified" gorm:"not null;default:false"`
	VerifiedAt      *time.Time `json:"verified_at"`
	User            User       `json:"-" gorm:"constraint:OnDelete:CASCADE"`
}
