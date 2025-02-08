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

	"github.com/vbvictor/ccv/pkg/complexity"

	"github.com/spf13/pflag"
	"golang.org/x/exp/maps"
)

// SortType represents the type of sorting to be performed on the results of git log.
type SortType = string

var (
	Changes   SortType = "changes"
	Additions SortType = "additions"
	Deletions SortType = "deletions"
	Commits   SortType = "commits"
)

const (
	DefaultTop = 10
	HashLength = 40
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
		return fmt.Errorf("failed to parse time: %w", err)
	}

	*d = Date{parsedTime}

	return nil
}

type ChurnOptions struct {
	SortBy       SortType
	Top          int
	Path         string
	ExcludePath  string
	Extensions   map[string]struct{}
	Since        Date
	Until        Date
	OutputFormat OutputType
}

var ChurnOpts = ChurnOptions{
	SortBy:       Changes,
	Top:          DefaultTop,
	Path:         "",
	ExcludePath:  "",
	Extensions:   nil,
	Since:        Date{},
	Until:        Date{},
	OutputFormat: Tabular,
}

func PrintRepoStats(repoPath string) error {
	churns, err := MostGitChurnFiles(repoPath)
	if err != nil {
		return fmt.Errorf("error getting churn metrics: %w", err)
	}

	return printStats(churns, os.Stdout, ChurnOpts)
}

func MostGitChurnFiles(repoPath string) ([]*complexity.ChurnChunk, error) {
	return ReadGitChurn(repoPath, ChurnOpts)
}

func ReadGitChurn(repoPath string, opts ChurnOptions) ([]*complexity.ChurnChunk, error) {
	cmd := buildGitCommand(repoPath, opts)

	output, err := executeGitCommand(cmd, repoPath)
	if err != nil {
		return nil, err
	}

	fileStats := make(map[string]*complexity.ChurnChunk)
	lines := strings.Split(string(output), "\n")

	processLines(lines, fileStats, opts)

	return sortAndLimit(maps.Values(fileStats), opts.SortBy, opts.Top), nil
}

func buildGitCommand(repoPath string, opts ChurnOptions) []string {
	cmd := []string{"git", "log", "--pretty=format:%H", "--numstat"}

	if !opts.Since.IsZero() {
		cmd = append(cmd, "--since="+opts.Since.String())
	}

	if !opts.Until.IsZero() {
		cmd = append(cmd, "--until="+opts.Until.String())
	}

	cmd = append(cmd, "--", repoPath)

	return cmd
}

func executeGitCommand(cmd []string, repoPath string) ([]byte, error) {
	gitCmd := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec // This command is built via buildGitCommand
	gitCmd.Dir = repoPath

	output, err := gitCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git command: %w", err)
	}

	return output, nil
}

func processLines(lines []string, fileStats map[string]*complexity.ChurnChunk, opts ChurnOptions) {
	currentCommit := ""
	modifiedInCommit := make(map[string]bool)

	for _, line := range lines {
		if line == "" {
			continue
		}

		if len(line) == HashLength {
			processCommit(currentCommit, modifiedInCommit, fileStats, opts)
			currentCommit = line
			modifiedInCommit = make(map[string]bool)
		} else {
			processFileLine(line, fileStats, modifiedInCommit, opts)
		}
	}
}

func processCommit(currentCommit string, modifiedInCommit map[string]bool, fileStats map[string]*complexity.ChurnChunk, opts ChurnOptions) {
	if currentCommit != "" && len(modifiedInCommit) > 0 {
		for filepath := range modifiedInCommit {
			if shouldSkipFile(filepath, opts) {
				continue
			}

			fileStats[filepath].Commits++
		}
	}
}

func processFileLine(line string, fileStats map[string]*complexity.ChurnChunk, modifiedInCommit map[string]bool, opts ChurnOptions) {
	parts := strings.Fields(line)
	if len(parts) == 3 && isNumeric(parts[0]) && isNumeric(parts[1]) {
		additions, _ := strconv.Atoi(parts[0])
		deletions, _ := strconv.Atoi(parts[1])

		path := parts[2]

		if shouldSkipFile(path, opts) {
			return
		}

		updateFileStats(fileStats, path, additions, deletions)

		modifiedInCommit[path] = true
	}
}

func updateFileStats(fileStats map[string]*complexity.ChurnChunk, path string, additions, deletions int) {
	if _, exists := fileStats[path]; !exists {
		fileStats[path] = &complexity.ChurnChunk{File: path}
	}

	fileStats[path].Added += additions
	fileStats[path].Removed += deletions
	fileStats[path].Churn += additions + deletions
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)

	return err == nil
}

func shouldSkipFile(file string, opts ChurnOptions) bool {
	if opts.ExcludePath != "" {
		if matched, _ := regexp.MatchString(opts.ExcludePath, file); matched {
			return true
		}
	}

	if opts.Extensions != nil {
		fileExt := filepath.Ext(file)

		if fileExt == "" {
			return false
		}

		if _, exists := opts.Extensions[fileExt[1:]]; !exists {
			return true
		}
	}

	return false
}

func sortAndLimit(result []*complexity.ChurnChunk, sortBy SortType, limit int) []*complexity.ChurnChunk {
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
