package random

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func Hex(n int) (string, error) {
	if n < 0 {
		return "", fmt.Errorf("random byte count must not be negative")
	}
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("read random bytes: %w", err)
	}
	return b, nil
}
