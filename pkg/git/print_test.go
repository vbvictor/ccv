package git

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/vbvictor/ccv/pkg/complexity"

	"github.com/stretchr/testify/assert"
)

func TestPrintTable(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*complexity.ChurnChunk
		opts     ChurnOptions
		expected []string
	}{
		{
			name: "single file churn",
			input: []*complexity.ChurnChunk{
				{
					File:    "main.go",
					Churn:   100,
					Added:   80,
					Removed: 20,
					Commits: 5,
				},
			},
			opts: ChurnOptions{
				Top:    1,
				SortBy: "changes",
			},
			expected: []string{
				"Top 1 most modified files by changes",
				"CHANGES", "ADDED", "DELETED", "COMMITS", "FILEPATH",
				"100", "80", "20", "5", "main.go",
			},
		},
		{
			name: "multiple files churn",
			input: []*complexity.ChurnChunk{
				{
					File:    "path/to/foo.go",
					Churn:   50,
					Added:   30,
					Removed: 20,
					Commits: 3,
				},
				{
					File:    "bar.go",
					Churn:   150,
					Added:   100,
					Removed: 50,
					Commits: 10,
				},
			},
			opts: ChurnOptions{
				Top:    2,
				SortBy: "commits",
			},
			expected: []string{
				"Top 2 most modified files by commits",
				"CHANGES", "ADDED", "DELETED", "COMMITS", "FILEPATH",
				"50", "30", "20", "3", "path/to/foo.go",
				"150", "100", "50", "10", "bar.go",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			printTable(tc.input, &buf, tc.opts)

			output := buf.String()
			for _, exp := range tc.expected {
				if !strings.Contains(output, exp) {
					t.Errorf("Expected output to contain %q, but it didn't.\nGot: %s", exp, output)
				}
			}
		})
	}
}

func TestPrintJSON(t *testing.T) {
	var buf bytes.Buffer

	results := []*complexity.ChurnChunk{
		{
			File:    "main.go",
			Churn:   10,
			Added:   5,
			Removed: 5,
			Commits: 2,
		},
	}

	since, _ := time.Parse("2006-01-02", "2024-01-01")
	until, _ := time.Parse("2006-01-02", "2024-01-31")

	opts := ChurnOptions{
		Top:         1,
		SortBy:      "churn",
		Path:        "src/",
		ExcludePath: "vendor/",
		Extensions:  map[string]struct{}{"go": {}},
		Since:       Date{since},
		Until:       Date{until},
	}

	printJSON(results, &buf, opts)

	expected := `{
  "metadata": {
    "totalFiles": 1,
    "sortBy": "churn",
    "filters": {
      "path": "src/",
      "excludePattern": "vendor/",
      "extensions": "go",
      "dateRange": {
        "since": "2024-01-01",
        "until": "2024-01-31"
      }
    }
  },
  "files": [
    {
      "path": "main.go",
      "changes": 10,
      "additions": 5,
      "deletions": 5,
      "commits": 2
    }
  ]
}
`
	assert.JSONEq(t, expected, buf.String())
}
