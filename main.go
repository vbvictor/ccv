package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vbvictor/ccv/pkg/complexity"

	"github.com/spf13/cobra"
	"github.com/vbvictor/ccv/pkg/git"
	"github.com/vbvictor/ccv/pkg/plot"
	"github.com/vbvictor/ccv/pkg/process"
)

// File to store the output graph
var (
	outputFile                   = ""
	ComplexityFuncThreshold uint = 5
)

func main() {
	cmdPlot := &cobra.Command{
		Use:   "plot [flags] <churn_file> <complexity_file>",
		Short: "Compare code complexity and churn metrics",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return plot.ValidateRiskThresholds()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			churnFile := args[0]
			complexityFile := args[1]

			if process.Verbose {
				fmt.Printf("Processing files:\n  Churn: %s\n  Complexity: %s\n", churnFile, complexityFile)
			}

			// Read churn data
			cf, err := os.Open(churnFile)
			if err != nil {
				return fmt.Errorf("error opening churn file: %w", err)
			}
			defer cf.Close()

			churns, err := complexity.ReadChurn(cf)
			if err != nil {
				return fmt.Errorf("error reading churn data: %w", err)
			}

			// Read complexity data
			xf, err := os.Open(complexityFile)
			if err != nil {
				return fmt.Errorf("Error opening complexity file: %w\n", err)
			}
			defer xf.Close()

			lizard, err := complexity.ReadLizardXML(xf)
			if err != nil {
				return fmt.Errorf("Error reading complexity data: %w\n", err)
			}

			files, err := complexity.ParseLizard(lizard)
			if err != nil {
				return fmt.Errorf("Error parsing complexity data: %w\n", err)
			}

			// Prepare plot data
			files = process.ApplyFilters(files, process.ComplexityFilter{MinComplexity: ComplexityFuncThreshold}.Filter)
			entries := process.PreparePlotData(files, churns)

			// Generate output
			switch plot.OutputFormat {
			case plot.Tabular:
				if err := plot.CreateTableChart(entries, os.Stdout); err != nil {
					return fmt.Errorf("error creating tabular chart: %w\n", err)
				}
			case plot.CSV:
				if err := plot.CreateTableChart(entries, os.Stdout); err != nil {
					return fmt.Errorf("error creating csv chart: %w\n", err)
				}
			case plot.Scatter:
				if err := plot.CreateScatterChart(entries, &plot.NoopMapper{}, outputFile); err != nil {
					return fmt.Errorf("error creating scatter chart: %w\n", err)
				}
			default:
				return fmt.Errorf("Invalid output format: %s\n", plot.OutputFormat)
			}

			if process.Verbose && plot.OutputFormat == plot.Scatter {
				fmt.Printf("Chart generated: %s\n", outputFile)
			}

			return nil
		},
	}

	flags := cmdPlot.PersistentFlags()
	flags.StringVarP(&outputFile, "output", "o", "complexity_churn.html", "Output file path")
	flags.BoolVarP(&process.Verbose, "verbose", "v", false, "Enable verbose output")
	flags.StringVarP(&process.Plot, "plot-type", "t", "commits", "Specify OY plot type: [commits, changes]")
	flags.UintVarP(&ComplexityFuncThreshold, "min-complexity", "m", 5, "Complexity threshold to delete functions with low complexity from the plot")
	flags.StringVarP(&plot.OutputFormat, "output-format", "f", "tabular", "Specify output format: [tabular, csv, scatter]")

	cmdChurn := &cobra.Command{
		Use:   "churn <repository>",
		Short: "Get churn metrics of a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoPath, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("error getting absolute path: %w", err)
			}

			if process.Verbose {
				fmt.Printf("Processing repository: %s\n", repoPath)
			}

			return git.PrintRepoStats(repoPath)
		},
	}

	flags = cmdChurn.PersistentFlags()
	flags.IntVar(&git.ChurnOpts.CommitCount, "commits", 0, "Number of commits to analyze")
	flags.StringVar(&git.ChurnOpts.SortBy, "sort", "changes", fmt.Sprintf("Sort by: %s, %s, %s, %s", git.Changes, git.Additions, git.Deletions, git.Commits))
	flags.IntVar(&git.ChurnOpts.Top, "top", 10, "Number of top files to display")
	flags.BoolVar(&process.Verbose, "verbose", false, "Show detailed progress")
	flags.StringVar(&git.ChurnOpts.ExcludePath, "exclude", "", "Exclude files matching regex pattern")
	flags.StringVar(&git.ChurnOpts.Extensions, "ext", "", "Only include files with extensions in comma-separated list. For example h,hpp,c,cpp")
	flags.Var(&git.ChurnOpts.Since, "since", "Start date for analysis (YYYY-MM-DD)")
	flags.Var(&git.ChurnOpts.Until, "until", "End date for analysis (YYYY-MM-DD)")
	flags.StringVar(&git.ChurnOpts.OutputFormat, "format", git.Tabular, fmt.Sprintf("Output format %v", git.OutputFormats))

	cmdChurn.Flag("since").DefValue = "none"
	cmdChurn.Flag("until").DefValue = "none"

	rootCmd := &cobra.Command{Use: "ccv"}
	rootCmd.AddCommand(cmdPlot, cmdChurn)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
