package models

import "github.com/google/uuid"

type PasswordCredential struct {
	Base
	AuthIdentityID uuid.UUID    `json:"-" gorm:"type:uuid;not null;uniqueIndex"`
	PasswordHash   string       `json:"-" gorm:"type:text;not null"`
	AuthIdentity   AuthIdentity `json:"-" gorm:"foreignKey:AuthIdentityID;constraint:OnDelete:CASCADE"`
}
