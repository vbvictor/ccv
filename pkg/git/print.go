package git

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/vbvictor/ccv/pkg/read"
)

// OutputType represents the type of output to be generated of churn subcommand
type OutputType = string

var (
	JSON          OutputType = "json"
	Tabular       OutputType = "tabular"
	OutputFormats            = []OutputType{JSON, Tabular}
)

func printStats(results []*read.ChurnChunk, out io.Writer, opts ChurnOptions) error {
	switch opts.OutputFormat {
	case JSON:
		printJSON(results, out, opts)
	case Tabular:
		printTable(results, out, opts)
	default:
		return fmt.Errorf("Invalid output format. Use one of the following: %v", OutputFormats)
	}

	return nil
}

func printTable(results []*read.ChurnChunk, out io.Writer, opts ChurnOptions) {
	fmt.Fprintf(out, "\nTop %d most modified files (by %s):\n", opts.Top, opts.SortBy)
	fmt.Fprintln(out, strings.Repeat("-", 100))
	fmt.Fprintf(out, "%-8s %-8s %-8s %-8s %s\n", "CHANGES", "ADDED", "DELETED", "COMMITS", "FILEPATH")
	fmt.Fprintln(out, strings.Repeat("-", 100))

	for _, chunk := range results {
		fmt.Fprintf(out, "%-8d %-8d %-8d %-8d %s\n",
			chunk.Churn,
			chunk.Added,
			chunk.Removed,
			chunk.Commits,
			chunk.File)
	}
}

func printJSON(results []*read.ChurnChunk, out io.Writer, opts ChurnOptions) {
	output := struct {
		Metadata struct {
			TotalFiles int    `json:"total_files"`
			SortBy     string `json:"sort_by"`
			Filters    struct {
				Path           string `json:"path"`
				ExcludePattern string `json:"exclude_pattern"`
				Extensions     string `json:"extensions"`
				DateRange      struct {
					Since string `json:"since"`
					Until string `json:"until"`
				} `json:"date_range"`
			} `json:"filters"`
		} `json:"metadata"`
		Files []*read.ChurnChunk `json:"files"`
	}{
		Files: results,
	}

	output.Metadata.TotalFiles = len(results)
	output.Metadata.SortBy = opts.SortBy
	output.Metadata.Filters.Path = opts.Path
	output.Metadata.Filters.ExcludePattern = opts.ExcludePath
	output.Metadata.Filters.Extensions = opts.Extensions
	output.Metadata.Filters.DateRange.Since = opts.Since.String()
	output.Metadata.Filters.DateRange.Until = opts.Until.String()

	json, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(out, "Error creating JSON output: %v\n", err)
		return
	}
	fmt.Fprintln(out, string(json))
}
