package users

import (
	"fmt"
	"net/mail"
	"strings"
)

func NormalizeEmail(raw string) (display, normalized string, err error) {
	email := strings.TrimSpace(raw)
	if email == "" {
		return "", "", fmt.Errorf("email is required")
	}
	address, err := mail.ParseAddress(email)
	if err != nil {
		return "", "", fmt.Errorf("invalid email: %w", err)
	}
	if address.Address != email {
		return "", "", fmt.Errorf("invalid email")
	}

	display = email
	normalized = strings.ToLower(email)
	return display, normalized, nil
}
