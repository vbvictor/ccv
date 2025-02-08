package git

import (
	"os/exec"
	"testing"

	"github.com/vbvictor/ccv/pkg/complexity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testData = []*complexity.ChurnChunk{
	{File: "file1.go", Churn: 100, Added: 60, Removed: 40, Commits: 5},
	{File: "file2.go", Churn: 200, Added: 150, Removed: 50, Commits: 3},
	{File: "file3.go", Churn: 150, Added: 70, Removed: 80, Commits: 8},
	{File: "file4.go", Churn: 80, Added: 30, Removed: 60, Commits: 2},
}

func TestSortAndLimitTypes(t *testing.T) {
	tests := []struct {
		name     string
		sortBy   SortType
		expected []string
	}{
		{
			name:     "sort by changes",
			sortBy:   Changes,
			expected: []string{"file2.go", "file3.go", "file1.go", "file4.go"},
		},
		{
			name:     "sort by additions",
			sortBy:   Additions,
			expected: []string{"file2.go", "file3.go", "file1.go", "file4.go"},
		},
		{
			name:     "sort by deletions",
			sortBy:   Deletions,
			expected: []string{"file3.go", "file4.go", "file2.go", "file1.go"},
		},
		{
			name:     "sort by commits",
			sortBy:   Commits,
			expected: []string{"file3.go", "file1.go", "file2.go", "file4.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortAndLimit(testData, tt.sortBy, 10)

			actual := extractFileNames(result)
			assert.Equal(t, tt.expected, actual)
			assertSorted(t, result, func(cc *complexity.ChurnChunk) any {
				switch tt.sortBy {
				case Changes:
					return cc.Churn
				case Additions:
					return cc.Added
				case Deletions:
					return cc.Removed
				case Commits:
					return cc.Commits
				default:
					return nil
				}
			})
		})
	}
}

func TestSortAndLimitLimits(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		expected []string
	}{
		{
			name:     "limit 2",
			limit:    2,
			expected: []string{"file3.go", "file1.go"},
		},
		{
			name:     "limit 0",
			limit:    0,
			expected: []string{},
		},
		{
			name:     "limit negative",
			limit:    -1,
			expected: []string{"file3.go", "file1.go", "file2.go", "file4.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortAndLimit(testData, Commits, tt.limit)

			actual := extractFileNames(result)
			assert.Equal(t, tt.expected, actual)
			assertSorted(t, result, func(cc *complexity.ChurnChunk) any { return cc.Commits })
		})
	}
}

func TestSortAndLimitFiles(t *testing.T) {
	tests := []struct {
		name     string
		input    []*complexity.ChurnChunk
		expected []string
	}{
		{
			name:     "empty input",
			input:    []*complexity.ChurnChunk{},
			expected: []string{},
		},
		{
			name: "single file",
			input: []*complexity.ChurnChunk{
				{File: "single.go", Commits: 1},
			},
			expected: []string{"single.go"},
		},
		{
			name: "multiple identical values",
			input: []*complexity.ChurnChunk{
				{File: "file1.go", Commits: 5},
				{File: "file2.go", Commits: 10},
				{File: "file3.go", Commits: 7},
			},
			expected: []string{"file2.go", "file3.go", "file1.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortAndLimit(tt.input, Commits, 10)

			actual := extractFileNames(result)
			assert.Equal(t, tt.expected, actual)
			assertSorted(t, result, func(cc *complexity.ChurnChunk) any { return cc.Commits })
		})
	}
}

func extractFileNames(chunks []*complexity.ChurnChunk) []string {
	names := make([]string, len(chunks))
	for i, chunk := range chunks {
		names[i] = chunk.File
	}

	return names
}

func assertSorted(t *testing.T, result []*complexity.ChurnChunk, ext func(*complexity.ChurnChunk) any) {
	for i := 1; i < len(result); i++ {
		assert.GreaterOrEqual(t, ext(result[i-1]), ext(result[i]))
	}
}

// TODO: add more data to bundle.
func TestReadChurn(t *testing.T) {
	tmpDir := t.TempDir()

	Unbundle(t, "../../test/bundles/churn-test.bundle", tmpDir)

	for _, tt := range []struct {
		name     string
		sortBy   SortType
		top      int
		expected []*complexity.ChurnChunk
	}{
		{
			name:   "sort by changes top 2",
			sortBy: Changes,
			top:    2,
			expected: []*complexity.ChurnChunk{
				{File: "main.cpp", Added: 15, Removed: 8, Churn: 23, Commits: 4},
				{File: "main.go", Added: 7, Removed: 0, Churn: 7, Commits: 1},
			},
		},
		{
			name:   "sort by additions",
			sortBy: Additions,
			top:    2,
			expected: []*complexity.ChurnChunk{
				{File: "main.cpp", Added: 15, Removed: 8, Churn: 23, Commits: 4},
				{File: "main.go", Added: 7, Removed: 0, Churn: 7, Commits: 1},
			},
		},
		{
			name:   "sort by deletions",
			sortBy: Deletions,
			top:    1,
			expected: []*complexity.ChurnChunk{
				{File: "main.cpp", Added: 15, Removed: 8, Churn: 23, Commits: 4},
			},
		},
		{
			name:   "sort by commits",
			sortBy: Commits,
			top:    1,
			expected: []*complexity.ChurnChunk{
				{File: "main.cpp", Added: 15, Removed: 8, Churn: 23, Commits: 4},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			results, err := ReadGitChurn(tmpDir, ChurnOptions{SortBy: tt.sortBy, Top: tt.top})
			require.NoError(t, err)
			assert.Len(t, results, len(tt.expected))

			for i, exp := range tt.expected {
				assert.Equal(t, exp.File, results[i].File)
				assert.Equal(t, exp.Added, results[i].Added)
				assert.Equal(t, exp.Removed, results[i].Removed)
				assert.Equal(t, exp.Churn, results[i].Churn)
				assert.Equal(t, exp.Commits, results[i].Commits)
			}
		})
	}
}

func Unbundle(t *testing.T, src, dst string) {
	t.Helper()

	cmd := exec.Command("git", "clone", src, dst)
	require.NoError(t, cmd.Run())
}

func TestShouldSkipFile(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		opts     ChurnOptions
		expected bool
	}{
		{
			name: "exclude pattern matches",
			file: "vendor/some/pkg/file.go",
			opts: ChurnOptions{
				ExcludePath: "vendor/.*",
			},
			expected: true,
		},
		{
			name: "exclude pattern does not match",
			file: "src/pkg/file.go",
			opts: ChurnOptions{
				ExcludePath: "vendor/.*",
			},
			expected: false,
		},
		{
			name: "extension matches allowed list",
			file: "main.go",
			opts: ChurnOptions{
				Extensions: map[string]struct{}{
					"go": {},
				},
			},
			expected: false,
		},
		{
			name: "extension not in allowed list",
			file: "script.py",
			opts: ChurnOptions{
				Extensions: map[string]struct{}{
					"go":  {},
					"cpp": {},
				},
			},
			expected: true,
		},
		{
			name:     "no filters applied",
			file:     "any/path/file.txt",
			opts:     ChurnOptions{},
			expected: false,
		},
		{
			name: "both filters applied - file matches both",
			file: "src/main.go",
			opts: ChurnOptions{
				ExcludePath: "test/.*",
				Extensions: map[string]struct{}{
					"go": {},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldSkipFile(tt.file, tt.opts)
			assert.Equal(t, tt.expected, result)
		})
	}
}
