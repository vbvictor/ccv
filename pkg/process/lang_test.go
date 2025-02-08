package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExtMap(t *testing.T) {
	tests := []struct {
		name     string
		langs    []LangName
		expected map[LangExt]struct{}
	}{
		{
			name:     "single language with one extension",
			langs:    []LangName{"go"},
			expected: map[LangExt]struct{}{"go": {}},
		},
		{
			name:     "multiple languages with multiple extensions",
			langs:    []LangName{"go", "py", "cpp", "hpp"},
			expected: map[LangExt]struct{}{"go": {}, "py": {}, "hpp": {}, "cpp": {}},
		},
		{
			name:     "empty language list",
			langs:    []LangName{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetExtMap(tt.langs)
			assert.Equal(t, tt.expected, got)
		})
	}
}
