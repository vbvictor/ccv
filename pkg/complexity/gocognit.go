package complexity

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/uudashr/gocognit"
)

func RunGocognit(repoPath string, opts Options) (FilesStat, error) {
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	var excludeRegex *regexp.Regexp
	if opts.Exclude != "" {
		excludeRegex, err = regexp.Compile(opts.Exclude)
		if err != nil {
			return nil, fmt.Errorf("invalid exclude pattern: %w", err)
		}
	}

	result := make(FilesStat, 0)
	fileMap := make(map[string][]FunctionStat)

	err = filepath.Walk(absRepoPath, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path: %w", err)
		}

		if excludeRegex != nil && excludeRegex.MatchString(path) {
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		stats := gocognit.ComplexityStats(file, fileSet, nil)
		functions := make([]FunctionStat, 0, len(stats))

		for _, stat := range stats {
			relPath, err := filepath.Rel(absRepoPath, stat.Pos.Filename)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}

			functions = append(functions, FunctionStat{
				File:       relPath,
				Package:    []string{stat.PkgName},
				Name:       stat.FuncName,
				Line:       stat.Pos.Line,
				Complexity: stat.Complexity,
			})
		}

		if len(functions) > 0 {
			relPath, err := filepath.Rel(absRepoPath, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path for file: %w", err)
			}

			fileMap[relPath] = functions
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk repository: %w", err)
	}

	// Convert map to FilesStat using relative paths
	for filePath, functions := range fileMap {
		result = append(result, FileStat{
			Path:      filePath,
			Functions: functions,
		})
	}

	return result, nil
}
