package process

import (
	"github.com/vbvictor/ccv/pkg/complexity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vbvictor/ccv/pkg/plot"
)

func TestComplexityFilter_Filter(t *testing.T) {
	files := complexity.FilesStat{
		&complexity.FileStat{
			Path: "file1.go",
			Functions: []complexity.FunctionStat{
				{Name: "func1", Compexity: 5},
				{Name: "func2", Compexity: 10},
				{Name: "func3", Compexity: 15},
			},
		},
		&complexity.FileStat{
			Path: "file2.go",
			Functions: []complexity.FunctionStat{
				{Name: "func4", Compexity: 3},
				{Name: "func5", Compexity: 7},
			},
		},
		&complexity.FileStat{
			Path: "file3.go",
			Functions: []complexity.FunctionStat{
				{Name: "func6", Compexity: 2},
			},
		},
	}

	tests := []struct {
		name          string
		minComplexity uint
		wantFiles     int
		wantFuncs     map[string]int
	}{
		{
			name:          "filter complexity >= 10",
			minComplexity: 10,
			wantFiles:     1,
			wantFuncs:     map[string]int{"file1.go": 2},
		},
		{
			name:          "filter complexity >= 5",
			minComplexity: 5,
			wantFiles:     2,
			wantFuncs:     map[string]int{"file1.go": 3, "file2.go": 1},
		},
		{
			name:          "filter complexity >= 20",
			minComplexity: 20,
			wantFiles:     0,
			wantFuncs:     map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := ComplexityFilter{MinComplexity: tt.minComplexity}
			got := filter.Filter(files)

			assert.Len(t, got, tt.wantFiles)

			for _, file := range got {
				wantFuncCount, exists := tt.wantFuncs[file.Path]
				assert.True(t, exists, "unexpected file in result: %s", file.Path)
				assert.Len(t, file.Functions, wantFuncCount)

				for _, fn := range file.Functions {
					assert.GreaterOrEqual(t, fn.Compexity, tt.minComplexity)
				}
			}
		})
	}
}

func TestApplyFilters(t *testing.T) {
	files := complexity.FilesStat{
		&complexity.FileStat{
			Path: "file1.go",
			Functions: []complexity.FunctionStat{
				{Name: "func1", Compexity: 5},
				{Name: "func2", Compexity: 10},
				{Name: "func3", Compexity: 15},
			},
		},
	}

	tests := []struct {
		name      string
		filters   []FilesFilterFunc
		wantFuncs int
	}{
		{
			name:      "no filters",
			filters:   []FilesFilterFunc{},
			wantFuncs: 3,
		},
		{
			name:      "single filter",
			filters:   []FilesFilterFunc{ComplexityFilter{MinComplexity: 7}.Filter},
			wantFuncs: 2,
		},
		{
			name:      "multiple filters",
			filters:   []FilesFilterFunc{ComplexityFilter{MinComplexity: 7}.Filter, ComplexityFilter{MinComplexity: 12}.Filter},
			wantFuncs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyFilters(files, tt.filters...)
			assert.Len(t, got[0].Functions, tt.wantFuncs)
		})
	}
}

func TestAvgComplexity(t *testing.T) {
	files := complexity.FilesStat{
		&complexity.FileStat{
			Path: "file1.go",
			Functions: []complexity.FunctionStat{
				{Name: "func1", Compexity: 5},
				{Name: "func2", Compexity: 10},
				{Name: "func3", Compexity: 15},
			},
		},
		&complexity.FileStat{
			Path: "file2.go",
			Functions: []complexity.FunctionStat{
				{Name: "func4", Compexity: 20},
				{Name: "func5", Compexity: 40},
			},
		},
		&complexity.FileStat{
			Path:      "empty.go",
			Functions: []complexity.FunctionStat{},
		},
	}

	got := avgComplexity(files)

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

func TestPreparePlotData(t *testing.T) {
	files := complexity.FilesStat{
		&complexity.FileStat{
			Path: "file1.go",
			Functions: []complexity.FunctionStat{
				{Name: "func1", Compexity: 5},
				{Name: "func2", Compexity: 10},
				{Name: "func3", Compexity: 15},
			},
		},
		&complexity.FileStat{
			Path: "file2.go",
			Functions: []complexity.FunctionStat{
				{Name: "func4", Compexity: 20},
				{Name: "func5", Compexity: 40},
			},
		},
		&complexity.FileStat{
			Path: "file3.go", // Will be skipped - no churn data
			Functions: []complexity.FunctionStat{
				{Name: "func6", Compexity: 25},
			},
		},
	}

	churns := []*complexity.ChurnChunk{
		{
			File:    "file1.go",
			Churn:   100,
			Added:   80,
			Removed: 20,
			Commits: 5,
		},
		{
			File:    "file2.go",
			Churn:   50,
			Added:   30,
			Removed: 20,
			Commits: 3,
		},
		{
			File:    "other.go", // Will be skipped - no complexity data
			Churn:   75,
			Added:   45,
			Removed: 30,
			Commits: 4,
		},
	}

	Plot = Changes
	got := PreparePlotData(files, churns)

	assert.Len(t, got, 2) // Only matching files should be included

	assert.Contains(t, got, plot.ScatterEntry{
		File:        "file1.go",
		ScatterData: plot.ScatterData{Complexity: 10, Churn: 100}, // (5 + 10 + 15) / 3
	})

	assert.Contains(t, got, plot.ScatterEntry{
		File:        "file2.go",
		ScatterData: plot.ScatterData{Complexity: 30, Churn: 50}, // (20 + 40) / 2
	})
}
