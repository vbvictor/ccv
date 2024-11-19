package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/vbvictor/ccv/pkg/read"
	"golang.org/x/exp/maps"
)

// SortType represents the type of sorting to be performed on the results of git log
type SortType = string

var (
	Changes   SortType = "changes"
	Additions SortType = "additions"
	Deletions SortType = "deletions"
	Commits   SortType = "commits"
)

var _ pflag.Value = (*Date)(nil)

type Date struct {
	time.Time
}

func (d *Date) Type() string {
	return "Date"
}

func (d *Date) String() string {
	return d.Format(time.DateOnly)
}

func (d *Date) Set(value string) error {
	parsedTime, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return err
	}
	*d = Date{parsedTime}
	return nil
}

type ChurnOptions struct {
	CommitCount  int
	SortBy       SortType
	Top          int
	Path         string
	ExcludePath  string
	Extensions   string
	Since        Date
	Until        Date
	OutputFormat OutputType
}

var ChurnOpts = ChurnOptions{
	CommitCount:  0,
	SortBy:       Changes,
	Top:          10,
	Path:         "",
	ExcludePath:  "",
	Extensions:   "",
	Since:        Date{},
	Until:        Date{},
	OutputFormat: Tabular,
}

func PrintRepoStats(repoPath string) error {
	churns, err := MostChurnFiles(repoPath)
	if err != nil {
		return fmt.Errorf("error getting churn metrics: %w", err)
	}

	return printStats(churns, os.Stdout, ChurnOpts)
}

func MostChurnFiles(repoPath string) ([]*read.ChurnChunk, error) {
	return ReadChurn(repoPath, ChurnOpts)
}

func ReadChurn(repoPath string, opts ChurnOptions) ([]*read.ChurnChunk, error) {
	cmd := []string{"git", "log", "--pretty=format:%H", "--numstat"}
	
	if opts.CommitCount > 0 {
		cmd = append(cmd, fmt.Sprintf("-n%d", opts.CommitCount))
	}

	if !opts.Since.IsZero() {
		cmd = append(cmd, fmt.Sprintf("--since=%s", opts.Since.String()))
	}

	if !opts.Until.IsZero() {
		cmd = append(cmd, fmt.Sprintf("--until=%s", opts.Until.String()))
	}

	cmd = append(cmd, "--", repoPath)

	gitCmd := exec.Command(cmd[0], cmd[1:]...)
	gitCmd.Dir = repoPath
	output, err := gitCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git command: %v", err)
	}

	fileStats := make(map[string]*read.ChurnChunk)
	lines := strings.Split(string(output), "\n")
	
	currentCommit := ""
	modifiedInCommit := make(map[string]bool)

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		if len(line) == 40 { // Commit hash
			if currentCommit != "" && len(modifiedInCommit) > 0 {
				for filepath := range modifiedInCommit {
					if shouldSkipFile(filepath, opts.ExcludePath, opts.Extensions) {
						continue
					}
					fileStats[filepath].Commits++
				}
			}
			currentCommit = line
			modifiedInCommit = make(map[string]bool)
		} else {
			parts := strings.Fields(line)
			if len(parts) == 3 && isNumeric(parts[0]) && isNumeric(parts[1]) {
				additions, _ := strconv.Atoi(parts[0])
				deletions, _ := strconv.Atoi(parts[1])
				filepath := parts[2]

				if shouldSkipFile(filepath, opts.ExcludePath, opts.Extensions) {
					continue
				}

				if _, exists := fileStats[filepath]; !exists {
					fileStats[filepath] = &read.ChurnChunk{File: filepath}
				}

				fileStats[filepath].Added += uint(additions)
				fileStats[filepath].Removed += uint(deletions)
				fileStats[filepath].Churn += uint(additions + deletions)
				modifiedInCommit[filepath] = true
			}
		}
	}

	result := maps.Values(fileStats)
	return sortAndLimit(result, opts.SortBy, opts.Top), nil
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func shouldSkipFile(file, excludePath, extensions string) bool {
	if excludePath != "" {
		if matched, _ := regexp.MatchString(excludePath, file); matched {
			return true
		}
	}

	if extensions != "" {
		ext := filepath.Ext(file)
		allowedExts := strings.Split(extensions, ",")
		found := false
		for _, allowedExt := range allowedExts {
			if "."+strings.TrimSpace(allowedExt) == ext {
				found = true
				break
			}
		}
		if !found {
			return true
		}
	}

	return false
}

func sortAndLimit(result []*read.ChurnChunk, sortBy SortType, limit int) []*read.ChurnChunk {
	less := func() func(i, j int) bool {
		switch sortBy {
		case Changes:
			return func(i, j int) bool { return result[i].Churn > result[j].Churn }
		case Additions:
			return func(i, j int) bool { return result[i].Added > result[j].Added }
		case Deletions:
			return func(i, j int) bool { return result[i].Removed > result[j].Removed }
		case Commits:
			return func(i, j int) bool { return result[i].Commits > result[j].Commits }
		default:
			return nil
		}
	}()

	sort.Slice(result, less)

	// Limit the number of results
	if limit >= 0 && len(result) > limit {
		result = result[:limit]
	}

	return result
}
