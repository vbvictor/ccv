package plot

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

func CreateTableChart(entries []ScatterEntry, out io.Writer) error {
	sort.Slice(entries, func(i, j int) bool {
		scoreI := entries[i].Complexity * float64(entries[i].Churn)
		scoreJ := entries[j].Complexity * float64(entries[j].Churn)
		return scoreI > scoreJ // Sort in descending order
	})

	fmt.Fprintln(out, "\nFiles ranked by risk score (Complexity * Churn):")
	fmt.Fprintln(out, strings.Repeat("-", 100))
	fmt.Fprintf(out, "%-12s %-12s %-12s %-s\n", "RISK SCORE", "COMPLEXITY", "CHURN", "FILEPATH")
	fmt.Fprintln(out, strings.Repeat("-", 100))

	for _, entry := range entries {
		riskScore := entry.Complexity * float64(entry.Churn)
		fmt.Fprintf(out, "%-12.2f %-12.2f %-12d %s\n",
			riskScore,
			entry.Complexity,
			entry.Churn,
			entry.File)
	}

	return nil
}
