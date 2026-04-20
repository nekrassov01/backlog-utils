package version

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	type globals struct {
		version  string
		revision string
	}
	type expected struct {
		value string
	}
	tests := []struct {
		name     string
		globals  globals
		expected expected
	}{
		{
			name: "basic",
			globals: globals{
				revision: "1234567",
			},
			expected: expected{
				value: fmt.Sprintf("%s (revision: 1234567)", version),
			},
		},
		{
			name: "no revision",
			globals: globals{
				revision: "",
			},
			expected: expected{
				value: version,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			revision = tt.globals.revision
			actual := Version()
			assert.Equal(t, tt.expected.value, actual)
		})
	}
}
