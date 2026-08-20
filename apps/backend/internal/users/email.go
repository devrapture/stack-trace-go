package users

import (
	"fmt"
	"strings"
)

func NormalizeEmail(raw string) (display, normalized string, err error) {
	email := strings.TrimSpace(raw)
	if email == "" {
		return "", "", fmt.Errorf("email is required")
	}

	atIndex := strings.LastIndex(email, "@")

	if atIndex <= 0 || atIndex == len(email)-1 {
		return "", "", fmt.Errorf("invalid email")
	}
	display = email
	normalized = strings.ToLower(email)
	return display, normalized, nil
}
