package apperrors

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrFailedToCreateUser is returned when a user fails to be created
	ErrFailedToCreateUser = errors.New("failed to create user")
)
