package complexity

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintTabular(t *testing.T) {
	testCases := []struct {
		name     string
		input    []FileComplexity
		expected []string
	}{
		{
			name: "single file complexity",
			input: []FileComplexity{
				{
					File:       "main.go",
					Complexity: 4,
				},
			},
			expected: []string{
				"main.go",
				"4",
			},
		},
		{
			name: "multiple files complexity",
			input: []FileComplexity{
				{
					File:       "path/to/foo.go",
					Complexity: 2,
				},
				{
					File:       "bar.go",
					Complexity: 6,
				},
			},
			expected: []string{
				"path/to/foo.go",
				"2",
				"bar.go",
				"6",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			PrintTabular(tc.input, &buf)

			output := buf.String()
			for _, exp := range tc.expected {
				if !strings.Contains(output, exp) {
					t.Errorf("Expected output to contain %q, but it didn't.\nGot: %s", exp, output)
				}
			}
		})
	}
}
