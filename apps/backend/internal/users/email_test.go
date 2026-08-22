package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantDisplay    string
		wantNormalized string
		wantErr        bool
	}{
		{
			name:           "normal email",
			input:          "alice@example.com",
			wantDisplay:    "alice@example.com",
			wantNormalized: "alice@example.com",
		},
		{
			name:           "mixed casing",
			input:          "Alice@Example.COM",
			wantDisplay:    "Alice@Example.COM",
			wantNormalized: "alice@example.com",
		},
		{
			name:           "surrounding spaces",
			input:          "  alice@example.com  ",
			wantDisplay:    "alice@example.com",
			wantNormalized: "alice@example.com",
		},
		{
			name:           "gmail dots remain",
			input:          "a.l.i.c.e@gmail.com",
			wantDisplay:    "a.l.i.c.e@gmail.com",
			wantNormalized: "a.l.i.c.e@gmail.com",
		},
		{
			name:           "plus tag remains",
			input:          "alice+work@gmail.com",
			wantDisplay:    "alice+work@gmail.com",
			wantNormalized: "alice+work@gmail.com",
		},
		{
			name:    "empty",
			input:   "",
			wantErr: true,
		},
		{
			name:    "only whitespace",
			input:   "   ",
			wantErr: true,
		},
		{
			name:    "missing at symbol",
			input:   "alice.example.com",
			wantErr: true,
		},
		{
			name:    "missing local part",
			input:   "@example.com",
			wantErr: true,
		},
		{
			name:    "missing domain",
			input:   "alice@",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			display, normalized, err := NormalizeEmail(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tt.wantDisplay, display)
			assert.Equal(t, tt.wantNormalized, normalized)
		})
	}
}
