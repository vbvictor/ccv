package complexity

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestReadChurn(t *testing.T) {
	jsonData := `{
        "files": [
            {
                "path": "src/file1.cpp",
                "changes": 150,
                "additions": 100,
                "deletions": 50,
                "commits": 10
            },
            {
                "path": "path/to/src/file2.cpp",
                "changes": 75,
                "additions": 45,
                "deletions": 30,
                "commits": 5
            }
        ]
    }`

	reader := strings.NewReader(jsonData)
	got, err := ReadChurn(reader)

	assert.NoError(t, err)
	assert.Len(t, got, 2)

	assert.Contains(t, got, &ChurnChunk{
		File:    "src/file1.cpp",
		Churn:   150,
		Added:   100,
		Removed: 50,
		Commits: 10,
	})

	assert.Contains(t, got, &ChurnChunk{
		File:    "path/to/src/file2.cpp",
		Churn:   75,
		Added:   45,
		Removed: 30,
		Commits: 5,
	})
}
