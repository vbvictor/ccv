package complexity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComplexityFilter(t *testing.T) {
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
				{Name: "func4", Complexity: 3},
				{Name: "func5", Complexity: 7},
			},
		},
		FileStat{
			Path: "file3.go",
			Functions: []FunctionStat{
				{Name: "func6", Complexity: 2},
			},
		},
	}

	tests := []struct {
		name          string
		minComplexity int
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
			filter := MinComplexityFilter{MinComplexity: tt.minComplexity}
			got := filter.Filter(files)

			assert.Len(t, got, tt.wantFiles)

			for _, file := range got {
				wantFuncCount, exists := tt.wantFuncs[file.Path]
				assert.True(t, exists, "unexpected file in result: %s", file.Path)
				assert.Len(t, file.Functions, wantFuncCount)

				for _, fn := range file.Functions {
					assert.GreaterOrEqual(t, fn.Complexity, tt.minComplexity)
				}
			}
		})
	}
}

func TestApplyFilters(t *testing.T) {
	files := FilesStat{
		FileStat{
			Path: "file1.go",
			Functions: []FunctionStat{
				{Name: "func1", Complexity: 5},
				{Name: "func2", Complexity: 10},
				{Name: "func3", Complexity: 15},
			},
		},
	}

	tests := []struct {
		name      string
		filters   []FilesFilterFunc
		wantFuncs []FunctionStat
	}{
		{
			name:    "no filters",
			filters: []FilesFilterFunc{},
			wantFuncs: []FunctionStat{
				{Name: "func1", Complexity: 5},
				{Name: "func2", Complexity: 10},
				{Name: "func3", Complexity: 15},
			},
		},
		{
			name:    "single filter",
			filters: []FilesFilterFunc{MinComplexityFilter{MinComplexity: 7}.Filter},
			wantFuncs: []FunctionStat{
				{Name: "func2", Complexity: 10},
				{Name: "func3", Complexity: 15},
			},
		},
		{
			name: "multiple filters",
			filters: []FilesFilterFunc{
				MinComplexityFilter{MinComplexity: 7}.Filter, MinComplexityFilter{MinComplexity: 12}.Filter,
			},
			wantFuncs: []FunctionStat{
				{Name: "func3", Complexity: 15},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyFilters(files, tt.filters...)
			assert.Equal(t, tt.wantFuncs, got[0].Functions)
		})
	}
}
