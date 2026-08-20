package models

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

type User struct {
	Base

	DisplayName string     `json:"display_name" gorm:"not null"`
	Status      UserStatus `json:"status" gorm:"default:active;check:chk_users_status, status IN ('active','suspended','deleted');not null"`

	Email          []UserEmail    `json:"-" gorm:"foreignKey:UserID"`
	AuthIdentities []AuthIdentity `json:"-" gorm:"foreignKey:UserID"`
}
