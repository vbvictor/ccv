package plot

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/vbvictor/ccv/pkg/complexity"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/assert"
)

func TestGroupByFile(t *testing.T) {
	tests := []struct {
		name    string
		entries []ScatterEntry
		want    []groupedEntry
	}{
		{
			name:    "empty entries returns empty result",
			entries: []ScatterEntry{},
			want:    []groupedEntry{},
		},
		{
			name: "single file entry returns single group",
			entries: []ScatterEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file1.go",
				},
			},
			want: []groupedEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					Files:       []string{"file1.go"},
				},
			},
		},
		{
			name: "multiple files with same metrics group together",
			entries: []ScatterEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file1.go",
				},
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file2.go",
				},
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file3.go",
				},
			},
			want: []groupedEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					Files:       []string{"file1.go", "file2.go", "file3.go"},
				},
			},
		},
		{
			name: "multiple files multiple metrics group",
			entries: []ScatterEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file1.go",
				},
				{
					ScatterData: ScatterData{Complexity: 9.0, Churn: 3},
					File:        "file2.go",
				},
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file3.go",
				},
				{
					ScatterData: ScatterData{Complexity: 7.0, Churn: 3},
					File:        "file4.go",
				},
				{
					ScatterData: ScatterData{Complexity: 9.0, Churn: 5},
					File:        "file5.go",
				},
			},
			want: []groupedEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					Files:       []string{"file1.go", "file3.go"},
				},
				{
					ScatterData: ScatterData{Complexity: 9.0, Churn: 3},
					Files:       []string{"file2.go"},
				},
				{
					ScatterData: ScatterData{Complexity: 7.0, Churn: 3},
					Files:       []string{"file4.go"},
				},
				{
					ScatterData: ScatterData{Complexity: 9.0, Churn: 5},
					Files:       []string{"file5.go"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := groupByFile(tt.entries)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

type stubMapper struct{}

func (m *stubMapper) Map(data ScatterData) Category {
	if data.Complexity >= 10.0 {
		return "critical"
	}
	if data.Complexity >= 5.0 {
		return "warning"
	}
	return "normal"
}

func (m *stubMapper) Style(_ Category) opts.ItemStyle {
	return opts.ItemStyle{}
}

func TestFormDataSeries(t *testing.T) {
	tests := []struct {
		name    string
		entries []ScatterEntry
		want    ScatterSeries
	}{
		{
			name:    "empty entries returns empty series",
			entries: []ScatterEntry{},
			want:    ScatterSeries{},
		},
		{
			name: "entries are grouped by same metrics and mapped to categories",
			entries: []ScatterEntry{
				{
					ScatterData: ScatterData{Complexity: 12.0, Churn: 5},
					File:        "critical1.go",
				},
				{
					ScatterData: ScatterData{Complexity: 12.0, Churn: 5},
					File:        "critical2.go",
				},
				{
					ScatterData: ScatterData{Complexity: 7.0, Churn: 3},
					File:        "warning.go",
				},
				{
					ScatterData: ScatterData{Complexity: 3.0, Churn: 1},
					File:        "normal1.go",
				},
				{
					ScatterData: ScatterData{Complexity: 1.0, Churn: 1},
					File:        "normal2.go",
				},
				{
					ScatterData: ScatterData{Complexity: 2.0, Churn: 1},
					File:        "normal3.go",
				},
			},
			want: ScatterSeries{
				"critical": []opts.ScatterData{
					{
						Value:      []interface{}{12.0, uint(5), "critical1.go<br/>critical2.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
				},
				"warning": []opts.ScatterData{
					{
						Value:      []interface{}{7.0, uint(3), "warning.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
				},
				"normal": []opts.ScatterData{
					{
						Value:      []interface{}{3.0, uint(1), "normal1.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
					{
						Value:      []interface{}{1.0, uint(1), "normal2.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
					{
						Value:      []interface{}{2.0, uint(1), "normal3.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
				},
			},
		},
	}

	mapper := &stubMapper{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formDataSeries(tt.entries, mapper)

			for category, series := range tt.want {
				assert.ElementsMatch(t, series, got[category])
			}
		})
	}
}

func TestPreparePlotData(t *testing.T) {
	files := complexity.FilesStat{
		complexity.FileStat{
			Path: "file1.go",
			Functions: []complexity.FunctionStat{
				{Name: "func1", Compexity: 5},
				{Name: "func2", Compexity: 10},
				{Name: "func3", Compexity: 15},
			},
		},
		complexity.FileStat{
			Path: "file2.go",
			Functions: []complexity.FunctionStat{
				{Name: "func4", Compexity: 20},
				{Name: "func5", Compexity: 40},
			},
		},
		complexity.FileStat{
			Path: "file3.go", // Will be skipped - no churn data
			Functions: []complexity.FunctionStat{
				{Name: "func6", Compexity: 25},
			},
		},
	}

	churns := []*complexity.ChurnChunk{
		{
			File:    "file1.go",
			Churn:   100,
			Added:   80,
			Removed: 20,
			Commits: 5,
		},
		{
			File:    "file2.go",
			Churn:   50,
			Added:   30,
			Removed: 20,
			Commits: 3,
		},
		{
			File:    "other.go", // Will be skipped - no complexity data
			Churn:   75,
			Added:   45,
			Removed: 30,
			Commits: 4,
		},
	}

	Plot = Changes
	got := PreparePlotData(files, churns)

	assert.Len(t, got, 2) // Only matching files should be included

	assert.Contains(t, got, ScatterEntry{
		File:        "file1.go",
		ScatterData: ScatterData{Complexity: 10, Churn: 100}, // (5 + 10 + 15) / 3
	})

	assert.Contains(t, got, ScatterEntry{
		File:        "file2.go",
		ScatterData: ScatterData{Complexity: 30, Churn: 50}, // (20 + 40) / 2
	})
}


func readCSVChartEntries(filepath string) ([]ScatterEntry, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	entries := make([]ScatterEntry, 0, len(records))
	for _, record := range records {
		complexity, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}

		churn, err := strconv.ParseUint(record[2], 10, 64)
		if err != nil {
			return nil, err
		}

		entry := ScatterEntry{
			File:        record[0],
			ScatterData: ScatterData{Complexity: complexity, Churn: uint(churn)},
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func createTestChart(t *testing.T, entries []ScatterEntry, outputPath string) { 
	err := CreateScatterChart(entries, NewRisksMapper(), outputPath)
	assert.NoError(t, err)

	_, err = os.Stat(outputPath)
	assert.NoError(t, err)
}

var CSVDataDir = "../../test/data/"
var OutputDir = "../../test/charts/"

func TestCreateScatterCharts(t *testing.T) {
	err := os.MkdirAll(OutputDir, 0755)
	assert.NoError(t, err)

	testCases := []struct {
		name     string
		csvFile  string
		outFile  string
	}{
		{
			name:    "200 different entries",
			csvFile: "plot_200.csv",
			outFile: "scatter-200.html",
		},
		{
			name:    "2000 different entries",
			csvFile: "plot_2000.csv",
			outFile: "scatter-2000.html",
		},
		{
			name:    "10 same entries",
			csvFile: "plot_10-same.csv",
			outFile: "scatter-10-same.html",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entries, err := readCSVChartEntries(CSVDataDir + tc.csvFile)
			assert.NoError(t, err)

			createTestChart(t, entries, OutputDir+tc.outFile)
		})
	}
}