package complexity

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunGocognit(t *testing.T) {
	tests := []struct {
		name          string
		directory     string
		expectedFiles int
		expectedFuncs map[string]TestFunction
	}{
		{
			name:          "Empty directory",
			directory:     "empty",
			expectedFiles: 0,
			expectedFuncs: map[string]TestFunction{},
		},
		{
			name:          "Nested directory structure",
			directory:     "nested",
			expectedFiles: 8,
			expectedFuncs: map[string]TestFunction{
				"BaseFunction":              {0, 3, "nested", "main.go"},
				"SimpleCondition":           {1, 7, "nested", "main.go"},
				"NestedIf":                  {3, 3, "level1", filepath.Join("level1", "file1.go")},
				"LoopWithCondition":         {6, 13, "level1", filepath.Join("level1", "file1.go")},
				"Func1":                     {6, 3, "level2", filepath.Join("level1", "level2", "file1.go")},
				"Func2":                     {3, 19, "level2", filepath.Join("level1", "level2", "file1.go")},
				"Func3":                     {6, 3, "level2", filepath.Join("level1", "level2", "file2.go")},
				"Func4":                     {10, 15, "level2", filepath.Join("level1", "level2", "file2.go")},
				"NestedLoopsWithConditions": {10, 3, "level2", filepath.Join("level1", "level2", "file3.go")},
				"SwitchWithLoops":           {9, 17, "level2", filepath.Join("level1", "level2", "file3.go")},
				"Func5":                     {15, 3, "level2", filepath.Join("level1", "level2", "morelevel2", "file1.go")},
				"Func6":                     {6, 19, "level2", filepath.Join("level1", "level2", "morelevel2", "file1.go")},
				"Func7":                     {14, 3, "level2", filepath.Join("level1", "level2", "morelevel2", "file2.go")},
				"Func8":                     {21, 24, "level2", filepath.Join("level1", "level2", "morelevel2", "file2.go")},
				"ComplexNestedStructure":    {21, 3, "level2", filepath.Join("level1", "level2", "morelevel2", "file3.go")},
				"MultipleControlFlows":      {21, 24, "level2", filepath.Join("level1", "level2", "morelevel2", "file3.go")},
			},
		},
		{
			name:          "Mixed complexity functions",
			directory:     "mixed",
			expectedFiles: 1,
			expectedFuncs: map[string]TestFunction{
				"simpleFunction":  {0, 3, "mixed", "main.go"},
				"complexFunction": {12, 7, "mixed", "main.go"},
			},
		},
		{
			name:          "Special cases",
			directory:     "special",
			expectedFiles: 1,
			expectedFuncs: map[string]TestFunction{
				"simpleFunction":  {0, 3, "special", "my main.go"},
				"complexFunction": {12, 7, "special", "my main.go"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPath := filepath.Join("..", "..", "test", "complexity", "gocode", tt.directory)

			result, err := RunGocognit(testPath, Options{})

			require.NoError(t, err)
			assert.Len(t, result, tt.expectedFiles)

			if tt.expectedFiles == 0 {
				return
			}

			foundFuncs := make(map[string]bool)

			for _, file := range result {
				for _, fn := range file.Functions {
					expected, exists := tt.expectedFuncs[fn.Name]
					assert.True(t, exists, "Function %s should exist", fn.Name)

					if exists {
						assert.Equal(t, expected.complexity, fn.Complexity,
							"Function %s should have complexity %d", fn.Name, expected.complexity)
						assert.Equal(t, expected.line, fn.Line,
							"Function %s should be on line %d", fn.Name, expected.line)
						assert.Equal(t, []string{expected.pkg}, fn.Package,
							"Function %s should be in package %s", fn.Name, expected.pkg)
						assert.Equal(t, expected.file, fn.File,
							"Function %s file path must exactly match %s, got %s",
							fn.Name, expected.file, fn.File)

						foundFuncs[fn.Name] = true
					}
				}
			}

			for funcName := range tt.expectedFuncs {
				assert.True(t, foundFuncs[funcName],
					"Expected function %s not found in results", funcName)
			}
		})
	}
}
