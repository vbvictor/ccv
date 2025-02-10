package plot

import (
	"fmt"
	"io"
	"sort"
)

func CreateCSVChart(entries []ScatterEntry, out io.Writer) error {
	sort.Slice(entries, func(i, j int) bool {
		scoreI := entries[i].Complexity * float64(entries[i].Churn)
		scoreJ := entries[j].Complexity * float64(entries[j].Churn)
		return scoreI > scoreJ
	})

	fmt.Fprintf(out, "RiskScore,Complexity,Churn,FilePath\n")

	for _, entry := range entries {
		riskScore := entry.Complexity * float64(entry.Churn)
		fmt.Fprintf(out, "%.2f,%.2f,%d,%s\n",
			riskScore,
			entry.Complexity,
			entry.Churn,
			entry.File)
	}

	return nil
}
