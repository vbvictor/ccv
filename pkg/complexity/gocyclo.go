package complexity

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/fzipp/gocyclo"
)

func RunGocyclo(repoPath string, opts Options) (FilesStat, error) {
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	}

	var excludeRegex *regexp.Regexp
	if opts.Exclude != "" {
		excludeRegex, err = regexp.Compile(opts.Exclude)
		if err != nil {
			return nil, fmt.Errorf("invalid exclude pattern: %w", err)
		}
	}

	paths := []string{absRepoPath}
	// Use gocyclo's built-in ignore functionality
	stats := gocyclo.Analyze(paths, excludeRegex)

	result := make(FilesStat, 0)
	fileMap := make(map[string][]FunctionStat)

	for _, stat := range stats {
		relPath, err := filepath.Rel(absRepoPath, stat.Pos.Filename)
		if err != nil {
			return nil, err
		}

		funcStat := FunctionStat{
			File:       relPath,
			Package:    []string{stat.PkgName},
			Name:       stat.FuncName,
			Line:       stat.Pos.Line,
			Complexity: stat.Complexity,
		}

		fileMap[relPath] = append(fileMap[relPath], funcStat)
	}

	for filePath, functions := range fileMap {
		result = append(result, FileStat{
			Path:      filePath,
			Functions: functions,
		})
	}

	return result, nil
}
