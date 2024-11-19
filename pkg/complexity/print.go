package complexity

import (
	"fmt"
	"io"
	"strings"
)

func PrintTabular(results FilesStat, out io.Writer) {
	fmt.Fprintln(out, "\nCode complexity analysis results:")
	fmt.Fprintln(out, strings.Repeat("-", 100))
	fmt.Fprintf(out, "%-50s %-15s %-15s %s\n", "FILEPATH", "FUNCTIONS", "AVG COMPLEX", "MAX COMPLEX")
	fmt.Fprintln(out, strings.Repeat("-", 100))

	for _, file := range results {
		avgComplexity := 0.0
		maxComplexity := uint(0)

		for _, fn := range file.Functions {
			avgComplexity += float64(fn.Compexity)
			if fn.Compexity > maxComplexity {
				maxComplexity = fn.Compexity
			}
		}

		if len(file.Functions) > 0 {
			avgComplexity /= float64(len(file.Functions))
		}

		fmt.Fprintf(out, "%-50s %-15d %-15.2f %d\n",
			file.Path,
			len(file.Functions),
			avgComplexity,
			maxComplexity)
	}
}
