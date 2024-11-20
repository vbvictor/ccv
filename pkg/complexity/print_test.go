package complexity

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintTabular(t *testing.T) {
	var buf bytes.Buffer

	results := FilesStat{
		&FileStat{
			Path: "file1.go",
			Functions: []FunctionStat{
				{Name: "func1", Compexity: 5},
				{Name: "func2", Compexity: 10},
			},
		},
		&FileStat{
			Path: "file2.go",
			Functions: []FunctionStat{
				{Name: "func3", Compexity: 3},
			},
		},
	}

	PrintTabular(results, &buf)

	expected := "\nCode complexity analysis results:\n" +
		"----------------------------------------------------------------------------------------------------\n" +
		"FILEPATH                                           FUNCTIONS       AVG COMPLEX     MAX COMPLEX\n" +
		"----------------------------------------------------------------------------------------------------\n" +
		"file1.go                                           2               7.50            10\n" +
		"file2.go                                           1               3.00            3\n"

	assert.Equal(t, expected, buf.String())
}

func TestPrintTabularEmpty(t *testing.T) {
	var buf bytes.Buffer

	results := FilesStat{
		&FileStat{
			Path:      "empty.go",
			Functions: []FunctionStat{},
		},
	}

	PrintTabular(results, &buf)

	expected := "\nCode complexity analysis results:\n" +
		"----------------------------------------------------------------------------------------------------\n" +
		"FILEPATH                                           FUNCTIONS       AVG COMPLEX     MAX COMPLEX\n" +
		"----------------------------------------------------------------------------------------------------\n" +
		"empty.go                                           0               0.00            0\n"

	assert.Equal(t, expected, buf.String())
}

func TestPrintTabularNoFiles(t *testing.T) {
	var buf bytes.Buffer

	results := FilesStat{}

	PrintTabular(results, &buf)

	expected := "\nCode complexity analysis results:\n" +
		"----------------------------------------------------------------------------------------------------\n" +
		"FILEPATH                                           FUNCTIONS       AVG COMPLEX     MAX COMPLEX\n" +
		"----------------------------------------------------------------------------------------------------\n"

	assert.Equal(t, expected, buf.String())
}