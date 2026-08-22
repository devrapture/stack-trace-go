package models

import "github.com/google/uuid"

type AuthIdentityKind string

const (
	AuthIdentityKindPassword AuthIdentityKind = "password"
	AuthIdentityKindGoogle   AuthIdentityKind = "google"
	AuthIdentityKindGithub   AuthIdentityKind = "github"
)

type AuthIdentity struct {
	Base

	UserID uuid.UUID        `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:uq_auth_identity_user_kind,priority:1"`
	Kind   AuthIdentityKind `json:"kind" gorm:"not null;uniqueIndex:uq_auth_identity_user_kind,priority:2;check:chk_auth_identity_kind,kind IN ('password','google','github')"`
	User   User             `json:"-" gorm:"constraint:OnDelete:CASCADE"`
}
