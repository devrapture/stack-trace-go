package password

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	apperrors "github.com/devrapture/stack-trace/internal/errors"
	"github.com/devrapture/stack-trace/internal/platform/random"
	"golang.org/x/crypto/argon2"
)

const (
	MinPasswordCharacters = 8
	MaxPasswordBytes      = 1024
)

type Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultParams = Params{
	Memory:      19 * 1024,
	Iterations:  2,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

func Validate(plainText string) error {
	if !utf8.ValidString(plainText) {
		return fmt.Errorf("password must contain valid UTF-8 characters")
	}

	if utf8.RuneCountInString(plainText) < MinPasswordCharacters {
		return fmt.Errorf("password must contain at least %d character", MinPasswordCharacters)
	}

	if len(plainText) > MaxPasswordBytes {
		return apperrors.ErrPasswordTooLong
	}
	return nil
}

func Hash(plainText string, params Params) (string, error) {
	if err := Validate(plainText); err != nil {
		return "", err
	}
	salt, err := random.Bytes(int(params.SaltLength))
	if err != nil {
		return "", fmt.Errorf("failed to generate password salt: %w", err)
	}
	hash := argon2.IDKey([]byte(plainText), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)
	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encoded, nil
}

func Verify(
	plaintext string,
	encoded string,
) (bool, error) {
	if len(plaintext) > MaxPasswordBytes {
		return false, apperrors.ErrPasswordTooLong
	}

	params, salt, storedHash, err := decode(encoded)
	if err != nil {
		return false, err
	}

	storedHashLen := len(storedHash)
	if storedHashLen > math.MaxUint32 {
		return false, apperrors.ErrInvalidHash
	}
	hashLength := uint32(storedHashLen)

	computedHash := argon2.IDKey(
		[]byte(plaintext),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		hashLength,
	)

	return subtle.ConstantTimeCompare(
		storedHash,
		computedHash,
	) == 1, nil
}

func NeedsRehash(
	encoded string,
	current Params,
) (bool, error) {
	params, _, _, err := decode(encoded)
	if err != nil {
		return false, err
	}

	return params.Memory < current.Memory ||
			params.Iterations < current.Iterations ||
			params.Parallelism < current.Parallelism,
		nil
}

func decode(
	encoded string,
) (Params, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")

	if len(parts) != 6 ||
		parts[1] != "argon2id" {
		return Params{}, nil, nil, apperrors.ErrInvalidHash
	}

	var version int

	if _, err := fmt.Sscanf(
		parts[2],
		"v=%d",
		&version,
	); err != nil {
		return Params{}, nil, nil, apperrors.ErrInvalidHash
	}

	if version != argon2.Version {
		return Params{},
			nil,
			nil,
			apperrors.ErrIncompatibleVersion
	}

	var params Params

	if _, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&params.Memory,
		&params.Iterations,
		&params.Parallelism,
	); err != nil {
		return Params{}, nil, nil, apperrors.ErrInvalidHash
	}

	if params.Memory == 0 ||
		params.Iterations == 0 ||
		params.Parallelism == 0 {
		return Params{}, nil, nil, apperrors.ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(
		parts[4],
	)
	if err != nil {
		return Params{},
			nil,
			nil,
			apperrors.ErrInvalidHash
	}

	hash, err := base64.RawStdEncoding.DecodeString(
		parts[5],
	)
	if err != nil {
		return Params{},
			nil,
			nil,
			apperrors.ErrInvalidHash
	}

	if len(salt) == 0 || len(hash) == 0 {
		return Params{},
			nil,
			nil,
			apperrors.ErrInvalidHash
	}

	return params, salt, hash, nil
}
