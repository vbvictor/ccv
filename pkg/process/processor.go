package process

import (
	"fmt"
	"github.com/vbvictor/ccv/pkg/complexity"

	"github.com/vbvictor/ccv/pkg/plot"
)

type ComplexityInputType int

const (
	Lizard ComplexityInputType = iota
	ClangTidy
)

type ChurnInputType int

const (
	ModifiedScript ChurnInputType = iota
)

var Version = "0.0.1"

var ComplexityInput ComplexityInputType = Lizard

var ChurnInput ChurnInputType = ModifiedScript

type PlotType = string

const (
	Commits PlotType = "commits"
	Changes PlotType = "changes"
)

var Plot = Commits

// FilesFilter Place where it is used
type FilesFilter interface {
	Filter(files complexity.FilesStat) complexity.FilesStat
}

type ComplexityFilter struct {
	MinComplexity uint
}

func (f ComplexityFilter) Filter(files complexity.FilesStat) complexity.FilesStat {
	result := make(complexity.FilesStat, 0, len(files))

	for _, file := range files {
		filteredFuncs := make([]complexity.FunctionStat, 0)
		for _, fn := range file.Functions {
			if fn.Compexity >= f.MinComplexity {
				filteredFuncs = append(filteredFuncs, fn)
			}
		}

		if len(filteredFuncs) > 0 {
			newFile := &complexity.FileStat{
				Path:      file.Path,
				Functions: filteredFuncs,
			}
			result = append(result, newFile)
		}
	}

	return result
}

type FilesFilterFunc func(files complexity.FilesStat) complexity.FilesStat

func ApplyFilters(files complexity.FilesStat, filters ...FilesFilterFunc) complexity.FilesStat {
	result := files

	for _, filter := range filters {
		result = filter(result)
	}

	return result
}

type FileComplexity struct {
	File       string
	Complexity float64
}

// Calculates average complexity bases on functions in file: sum(funcComplexity) / funcCount
func avgComplexity(files complexity.FilesStat) []FileComplexity {
	result := make([]FileComplexity, 0, len(files))

	for _, file := range files {
		if len(file.Functions) == 0 {
			continue
		}

		var totalComplexity float64 = 0
		for _, fn := range file.Functions {
			totalComplexity += float64(fn.Compexity)
		}

		complexity := totalComplexity / float64(len(file.Functions))
		if Verbose {
			fmt.Printf("File: %s, Complexity: %f\n", file.Path, complexity)
		}

		result = append(result, FileComplexity{
			File:       file.Path,
			Complexity: complexity,
		})
	}

	return result
}

// Skip file if it is not found in chunk or files, first goes over all churns
// Matches based on filename
func PreparePlotData(files complexity.FilesStat, churns []*complexity.ChurnChunk) []plot.ScatterEntry {
	result := make([]plot.ScatterEntry, 0)

	// Calculate average complexity for each file
	fileComplexities := avgComplexity(files)

	// Create map for quick churn lookup
	churnMap := make(map[string]*complexity.ChurnChunk)
	for _, churn := range churns {
		churnMap[churn.File] = churn
	}

	// Match files with churns and create chart entries
	for _, fc := range fileComplexities {
		churn, exists := churnMap[fc.File]

		if !exists {
			continue
		}

		entry := plot.ScatterEntry{
			File:        fc.File,
			ScatterData: plot.ScatterData{Complexity: fc.Complexity, Churn: 0},
		}

		switch Plot {
		case Commits:
			entry.Churn = churn.Commits
		case Changes:
			entry.Churn = churn.Churn
		default:
			panic("Unknown plot type")
		}

		result = append(result, entry)
	}

	return result
}
