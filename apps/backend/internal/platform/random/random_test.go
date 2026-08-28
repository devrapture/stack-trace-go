package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHex(t *testing.T) {
	result, err := Hex(8)
	require.NoError(t, err)
	assert.Len(t, result, 16)

}

func TestBytes(t *testing.T) {
	result, err := Bytes(32)
	require.NoError(t, err)
	assert.Len(t, result, 32)
}

func TestHex_GeneratesDifferentValues(t *testing.T) {
	first, err := Hex(16)
	require.NoError(t, err)

	second, err := Hex(16)
	require.NoError(t, err)

	assert.NotEqual(t, first, second)
}

func TestBytes_GeneratesDifferentValues(t *testing.T) {
	first, err := Bytes(16)
	require.NoError(t, err)

	second, err := Bytes(16)
	require.NoError(t, err)

	assert.NotEqual(t, first, second)
}
