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

func main() {
	var outputFile = ""

	cmdPlot := &cobra.Command{
		Use:   "plot [flags] <repository>",
		Short: "Compare code complexity and churn metrics",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return plot.ValidateRiskThresholds()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repoPath, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("error getting absolute path: %w", err)
			}

			if process.Verbose {
				fmt.Printf("Processing repository: %s\n", repoPath)
			}

			churns, err := git.MostGitChurnFiles(repoPath)
			if err != nil {
				return fmt.Errorf("error reading churn data: %w", err)
			}

			if err := complexity.CheckLizardExecutable(); err != nil {
				return err;
			}

			fileStat, err := complexity.RunLizardCmd(repoPath, complexity.Opts)
			if err != nil {
				return fmt.Errorf("error running lizard command: %w", err)
			}

			fileStat = complexity.ApplyFilters(fileStat, complexity.ComplexityFilter{MinComplexity: 5}.Filter)
			plotEntries := plot.PreparePlotData(fileStat, churns)

			if err := plot.CreateScatterChart(plotEntries, plot.NewRisksMapper(), outputFile); err != nil {
				return fmt.Errorf("error creating chart: %w\n", err)
			}

			if process.Verbose {
				fmt.Printf("Chart generated: %s\n", outputFile)
			}

			return nil
		},
	}

	flags := cmdPlot.PersistentFlags()
	flags.StringVarP(&outputFile, "output", "o", "complexity_churn.html", "Output file path")
	flags.BoolVarP(&process.Verbose, "verbose", "v", false, "Enable verbose output")
	flags.StringVarP(&plot.Plot, "plot-type", "t", "changes", "Specify OY plot type")
	flags.IntVar(&git.ChurnOpts.Top, "top", 10, "Number of top files to display")
	flags.StringVar(&git.ChurnOpts.ExcludePath, "exclude", "", "Exclude files matching regex pattern")
	flags.StringSliceVar(&process.Languages, "lang", nil, "Only include files with given languages in comma-separated list. For example cpp,java")
	flags.Var(&git.ChurnOpts.Since, "since", "Start date for analysis (YYYY-MM-DD)")
	flags.Var(&git.ChurnOpts.Until, "until", "End date for analysis (YYYY-MM-DD)")
	flags.IntVar(&complexity.Opts.Threads, "t", 1, "Number of threads to run")

	/*
		  flags.UintVar(&plot.VeryLowRisk, "very-low-risk", 10, "Very Low Risk threshold")
			flags.UintVar(&plot.LowRisk, "low-risk", 15, "Low Risk threshold")
			flags.UintVar(&plot.MediumRisk, "medium-risk", 20, "Medium Risk threshold")
			flags.UintVar(&plot.HighRisk, "high-risk", 25, "High Risk threshold")
			flags.UintVar(&plot.VeryHighRisk, "very-high-risk", 30, "Very High Risk threshold")
			flags.UintVar(&plot.CriticalRisk, "critical-risk", 35, "Critical Risk threshold")
	*/

	var cmdChurn = &cobra.Command{
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

			git.ChurnOpts.Extensions = process.GetExtMap(process.Languages)

			return git.PrintRepoStats(repoPath)
		},
	}

	flags = cmdChurn.PersistentFlags()
	flags.StringVar(&git.ChurnOpts.SortBy, "sort", "changes", fmt.Sprintf("Sort by: %s, %s, %s, %s", git.Changes, git.Additions, git.Deletions, git.Commits))
	flags.IntVar(&git.ChurnOpts.Top, "top", 10, "Number of top files to display")
	flags.BoolVar(&process.Verbose, "verbose", false, "Show detailed progress")
	flags.StringVar(&git.ChurnOpts.ExcludePath, "exclude", "", "Exclude files matching regex pattern")
	flags.StringSliceVar(&process.Languages, "lang", nil, "Only include files with given languages in comma-separated list. For example cpp,java")
	flags.Var(&git.ChurnOpts.Since, "since", "Start date for analysis (YYYY-MM-DD)")
	flags.Var(&git.ChurnOpts.Until, "until", "End date for analysis (YYYY-MM-DD)")
	flags.StringVar(&git.ChurnOpts.OutputFormat, "format", git.Tabular, fmt.Sprintf("Output format %v", git.OutputFormats))

	cmdChurn.Flag("since").DefValue = "none"
	cmdChurn.Flag("until").DefValue = "none"

	var cmdComplexity = &cobra.Command{
		Use:   "complexity <path>",
		Short: "Get complexity metrics of a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("error getting absolute path: %w", err)
			}

			if process.Verbose {
				fmt.Printf("Processing repository: %s\n", path)
			}

			if err := complexity.CheckLizardExecutable(); err != nil {
				return err;
			}

			fileStat, err := complexity.RunLizardCmd(path, complexity.Opts)
			if err != nil {
				return fmt.Errorf("error running lizard command: %w", err)
			}

			avgComplexity := complexity.AvgComplexity(fileStat)

			complexity.PrintTabular(avgComplexity, os.Stdout)

			return nil
		},
	}

	flags = cmdComplexity.PersistentFlags()
	flags.StringVar(&complexity.Opts.Languages, "lang", "", "Only include files with languages in comma-separated list. For example cpp,python")
	flags.IntVar(&complexity.Opts.Threads, "t", 1, "Number of threads to run")

	var rootCmd = &cobra.Command{Use: "ccv"}
	rootCmd.AddCommand(cmdPlot, cmdChurn, cmdComplexity)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
