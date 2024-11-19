package plot

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readCSVToChartEntries(filepath string) ([]ScatterEntry, error) {
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
			entries, err := readCSVToChartEntries(CSVDataDir + tc.csvFile)
			assert.NoError(t, err)

			createTestChart(t, entries, OutputDir+tc.outFile)
		})
	}
}