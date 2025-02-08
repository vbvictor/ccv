package complexity

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLizard(t *testing.T) {
	lizard := &lizardXML{
		Measures: []lizardMeasure{
			{
				Type: "File",
				Items: []lizardItem{
					{
						Name:   "src/file1.cpp",
						Values: []int{1, 2, 3, 4},
					},
				},
			},
			{
				Type: "Function",
				Items: []lizardItem{
					{
						Name:   "pkg::some_func(...) at src/file1.cpp:1",
						Values: []int{1, 2, 3},
					},
					{
						Name:   "pkg::sub_pkg::SomeFunc(...) at src/file1.cpp:3",
						Values: []int{4, 5, 6},
					},
				},
			},
		},
	}

	got, err := parseLizard(lizard)
	require.NoError(t, err)
	assert.Len(t, got, 1) // One file with two functions

	file := got[0]
	assert.Equal(t, "src/file1.cpp", file.Path)
	assert.Len(t, file.Functions, 2)

	// Check functions using Contains
	assert.Contains(t, file.Functions, FunctionStat{
		Name:       "some_func",
		Package:    []string{"pkg"},
		Line:       1,
		Length:     2,
		File:       "src/file1.cpp",
		Complexity: 3,
	})

	assert.Contains(t, file.Functions, FunctionStat{
		Name:       "SomeFunc",
		Package:    []string{"pkg", "sub_pkg"},
		Line:       3,
		Length:     5,
		File:       "src/file1.cpp",
		Complexity: 6,
	})
}

func TestParseItem(t *testing.T) {
	tests := []struct {
		name    string
		input   lizardItem
		want    FunctionStat
		wantErr bool
	}{
		{
			name: "valid function with package",
			input: lizardItem{
				Name:   "pkg::subpkg::func_name at path/to/file.go:123",
				Values: []int{0, 3, 5},
			},
			want: FunctionStat{
				Name:       "func_name",
				Package:    []string{"pkg", "subpkg"},
				Line:       123,
				Length:     3,
				File:       "path/to/file.go",
				Complexity: 5,
			},
			wantErr: false,
		},
		{
			name: "valid function with parameters",
			input: lizardItem{
				Name:   "pkg::func_name(param1, param2) at path/file.go:456",
				Values: []int{0, 2, 3},
			},
			want: FunctionStat{
				Name:       "func_name",
				Package:    []string{"pkg"},
				Line:       456,
				Length:     2,
				File:       "path/file.go",
				Complexity: 3,
			},
			wantErr: false,
		},
		{
			name: "function with many packages",
			input: lizardItem{
				Name:   "pkg1::pkg2::pkg3::pkg4::pkg5::func_name(param1, param2) at path/file.go:456",
				Values: []int{0, 1, 4},
			},
			want: FunctionStat{
				Name:       "func_name",
				Package:    []string{"pkg1", "pkg2", "pkg3", "pkg4", "pkg5"},
				Line:       456,
				Length:     1,
				File:       "path/file.go",
				Complexity: 4,
			},
			wantErr: false,
		},
		{
			name: "function without package",
			input: lizardItem{
				Name:   "standalone_func at path/file.go:789",
				Values: []int{0, 0, 1},
			},
			want: FunctionStat{
				Name:       "standalone_func",
				Package:    []string{},
				Line:       789,
				Length:     0,
				File:       "path/file.go",
				Complexity: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid format",
			input: lizardItem{
				Name:   "invalid_format",
				Values: []int{1, 1, 1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseItem(tt.input)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err, "got unexpected error")
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Package, got.Package)
			assert.Equal(t, tt.want.Line, got.Line)
			assert.Equal(t, tt.want.File, got.File)
			assert.Equal(t, tt.want.Complexity, got.Complexity)
		})
	}
}

func TestReadLizardXML(t *testing.T) {
	xmlData := `<?xml version="1.0" ?>
<cppncss>
    <measure type="Function">
        <labels>
            <label>Nr.</label>
            <label>NCSS</label>
            <label>CCN</label>
        </labels>
        <item name="pkg::some_func(...) at src/file1.cpp:1">
            <value>11</value>
            <value>111</value>
            <value>1111</value>
        </item>
        <item name="pkg::sub_pkg::SomeFunc(...) at path/to/src/file2.cpp:3">
            <value>22</value>
            <value>222</value>
            <value>2222</value>
        </item>
        <average label="NCSS" value="16.2"/>
        <average label="CCN" value="4.4"/>
    </measure>
    <measure type="File">
        <labels>
            <label>Nr.</label>
            <label>NCSS</label>
            <label>CCN</label>
            <label>Functions</label>
        </labels>
        <item name="src/file1.cpp">
            <value>1</value>
            <value>2</value>
            <value>3</value>
            <value>4</value>
        </item>
        <item name="path/to/src/file2.cpp">
            <value>5</value>
            <value>6</value>
            <value>7</value>
            <value>8</value>
        </item>
        <average label="NCSS" value="1"/>
        <average label="CCN" value="2"/>
        <average label="Functions" value="3"/>
        <sum label="NCSS" value="4"/>
        <sum label="CCN" value="5"/>
        <sum label="Functions" value="6"/>
    </measure>
</cppncss>`

	reader := strings.NewReader(xmlData)
	got, err := readLizardXML(reader)

	require.NoError(t, err)
	assert.NotNil(t, got)

	// Test Functions section
	idx := slices.IndexFunc(got.Measures, func(l lizardMeasure) bool {
		return l.Type == "Function"
	})
	assert.Greater(t, idx, -1)

	assert.Nil(t, got.Measures[idx].Sums)

	assert.Len(t, got.Measures[idx].Items, 2)
	assert.Contains(t, got.Measures[idx].Items, lizardItem{
		Name:   "pkg::some_func(...) at src/file1.cpp:1",
		Values: []int{11, 111, 1111},
	})
	assert.Contains(t, got.Measures[idx].Items, lizardItem{
		Name:   "pkg::sub_pkg::SomeFunc(...) at path/to/src/file2.cpp:3",
		Values: []int{22, 222, 2222},
	})

	// Test Files section
	idx = slices.IndexFunc(got.Measures, func(l lizardMeasure) bool {
		return l.Type == "File"
	})
	assert.Greater(t, idx, -1)

	assert.Len(t, got.Measures[idx].Items, 2)
	assert.Contains(t, got.Measures[idx].Items, lizardItem{
		Name:   "src/file1.cpp",
		Values: []int{1, 2, 3, 4},
	})
	assert.Contains(t, got.Measures[idx].Items, lizardItem{
		Name:   "path/to/src/file2.cpp",
		Values: []int{5, 6, 7, 8},
	})
}

func skipIfLizardNotPresent(t *testing.T) {
	t.Helper()

	if os.Getenv("ENABLE_LIZARD_EXECUTABLE") == "" {
		t.Skip("Skipping running lizard complexity analysis - set ENABLE_LIZARD_EXECUTABLE to enable")
	}
}

func TestRunLizardCmd(t *testing.T) {
	skipIfLizardNotPresent(t)

	for _, tt := range []struct {
		name          string
		language      string
		expectedFuncs map[string]uint
		file          string
	}{
		{
			name:     "calculate complexity from cpp files",
			language: "cpp",
			expectedFuncs: map[string]uint{
				"processNumber":           4,
				"validateAndProcessInput": 6,
			},
			file: "main.cpp",
		},
		{
			name:     "calculate complexity from go files",
			language: "python",
			expectedFuncs: map[string]uint{
				"calculate_grade":   4,
				"is_valid_password": 8,
			},
			file: "main.py",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			results, err := RunLizardCmd("../../test/complexity/lizard", Options{
				Exts:    tt.language,
				Threads: 1,
			})

			require.NoError(t, err)
			assert.NotEmpty(t, results)
			assert.Len(t, results, 1)
			assert.NotNil(t, results[0], "%s should be analyzed", tt.file)

			for _, fn := range results[0].Functions {
				expectedComplexity, exists := tt.expectedFuncs[fn.Name]
				if exists {
					assert.Equal(t, expectedComplexity, fn.Complexity,
						"Function %s should have complexity %d", fn.Name, expectedComplexity)
				}
			}
		})
	}
}

func TestCheckLizardExecutable(t *testing.T) {
	skipIfLizardNotPresent(t)

	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	tempDir := t.TempDir()
	fakeLizardPath := filepath.Join(tempDir, "lizard")

	err := os.WriteFile(fakeLizardPath, []byte(""), 0o600)
	if err != nil {
		t.Fatalf("failed to create fake lizard executable: %v", err)
	}

	tests := []struct {
		name      string
		pathValue string
		wantErr   bool
	}{
		{
			name:      "lizard exists in PATH",
			pathValue: tempDir + ":" + originalPath,
			wantErr:   false,
		},
		{
			name:      "lizard not in PATH",
			pathValue: "/nonexistent/path",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PATH", tt.pathValue)

			err := CheckLizardExecutable()
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, "lizard executable not found", err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
