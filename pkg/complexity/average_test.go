package complexity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvgComplexity(t *testing.T) {
	files := FilesStat{
		FileStat{
			Path: "file1.go",
			Functions: []FunctionStat{
				{Name: "func1", Complexity: 5},
				{Name: "func2", Complexity: 10},
				{Name: "func3", Complexity: 15},
			},
		},
		FileStat{
			Path: "file2.go",
			Functions: []FunctionStat{
				{Name: "func4", Complexity: 20},
				{Name: "func5", Complexity: 40},
			},
		},
		FileStat{
			Path:      "empty.go",
			Functions: []FunctionStat{},
		},
	}

	got := AvgComplexity(files)

	assert.Len(t, got, 2) // empty.go should be skipped

	assert.Contains(t, got, FileComplexity{
		File:       "file1.go",
		Complexity: 10, // (5 + 10 + 15) / 3
	})

	assert.Contains(t, got, FileComplexity{
		File:       "file2.go",
		Complexity: 30, // (20 + 40) / 2
	})
}
