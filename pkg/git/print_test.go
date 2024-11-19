package git

import (
	"bytes"
	"github.com/vbvictor/ccv/pkg/complexity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrintTable(t *testing.T) {
	var buf bytes.Buffer

	results := []*complexity.ChurnChunk{
		{
			File:    "main.go",
			Churn:   20,
			Added:   15,
			Removed: 5,
			Commits: 3,
		},
		{
			File:    "test.go",
			Churn:   10,
			Added:   5,
			Removed: 5,
			Commits: 2,
		},
	}

	opts := ChurnOptions{
		Top:    2,
		SortBy: "churn",
	}

	printTable(results, &buf, opts)

	expected := "\nTop 2 most modified files (by churn):\n" +
		"----------------------------------------------------------------------------------------------------\n" +
		"CHANGES  ADDED    DELETED  COMMITS  FILEPATH\n" +
		"----------------------------------------------------------------------------------------------------\n" +
		"20       15       5        3        main.go\n" +
		"10       5        5        2        test.go\n"

	assert.Equal(t, expected, buf.String())
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
		Extensions:  ".go,.ts",
		Since:       Date{since},
		Until:       Date{until},
	}

	printJSON(results, &buf, opts)

	expected := `{
  "metadata": {
    "total_files": 1,
    "sort_by": "churn",
    "filters": {
      "path": "src/",
      "exclude_pattern": "vendor/",
      "extensions": ".go,.ts",
      "date_range": {
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
