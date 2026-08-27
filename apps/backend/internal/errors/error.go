package apperrors

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrFailedToCreateUser is returned when a user fails to be created
	ErrFailedToCreateUser = errors.New("failed to create user")
	// ErrPasswordTooShort is returned when a password is too long
	ErrPasswordTooLong = errors.New("password is too long")
	// ErrInvalidHash is returned when a password hash is invalid
	ErrInvalidHash = errors.New("invalid password hash")
	// ErrIncompatibleVersion is returned when a password hash is incompatible with the current version
	ErrIncompatibleVersion = errors.New("incompatible argon2 version")
)
