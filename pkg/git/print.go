package git

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/bndr/gotabulate"
	"github.com/vbvictor/ccv/pkg/complexity"
	"golang.org/x/exp/maps"
)

// OutputType represents the type of output to be generated of churn subcommand.
type OutputType = string

var (
	JSON          OutputType = "json"
	Tabular       OutputType = "tabular"
	OutputFormats            = []OutputType{JSON, Tabular}
)

func printStats(results []*complexity.ChurnChunk, out io.Writer, opts ChurnOptions) error {
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

func printTable(results []*complexity.ChurnChunk, out io.Writer, opts ChurnOptions) {
	fmt.Fprintf(out, "\nTop %d most modified files by %s:\n", opts.Top, opts.SortBy)

	data := make([][]interface{}, len(results))

	for i, result := range results {
		data[i] = []interface{}{result.Churn, result.Added, result.Removed, result.Commits, result.File}
	}

	table := gotabulate.Create(data)
	table.SetHeaders([]string{"CHANGES", "ADDED", "DELETED", "COMMITS", "FILEPATH"})
	table.SetAlign("left")

	_, _ = io.WriteString(out, table.Render("grid"))
}

func printJSON(results []*complexity.ChurnChunk, out io.Writer, opts ChurnOptions) {
	output := struct {
		Metadata struct {
			TotalFiles int    `json:"totalFiles"`
			SortBy     string `json:"sortBy"`
			Filters    struct {
				Path           string `json:"path"`
				ExcludePattern string `json:"excludePattern"`
				Extensions     string `json:"extensions"`
				DateRange      struct {
					Since string `json:"since"`
					Until string `json:"until"`
				} `json:"dateRange"`
			} `json:"filters"`
		} `json:"metadata"`
		Files []*complexity.ChurnChunk `json:"files"`
	}{
		Files: results,
	}

	output.Metadata.TotalFiles = len(results)
	output.Metadata.SortBy = opts.SortBy
	output.Metadata.Filters.Path = opts.Path
	output.Metadata.Filters.ExcludePattern = opts.ExcludePath
	output.Metadata.Filters.Extensions = strings.Join(maps.Keys(opts.Extensions), ",")
	output.Metadata.Filters.DateRange.Since = opts.Since.String()
	output.Metadata.Filters.DateRange.Until = opts.Until.String()

	result, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(out, "Error creating JSON output: %v\n", err)

		return
	}

	fmt.Fprintln(out, string(result))
}
