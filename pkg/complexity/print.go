package complexity

import (
	"io"

	"github.com/bndr/gotabulate"
)

func PrintTabular(results []FileComplexity, out io.Writer) {
	_, _ = io.WriteString(out, "\nCode complexity analysis results:\n")

	data := make([][]interface{}, len(results))
	for i, result := range results {
		data[i] = []interface{}{result.File, result.Complexity}
	}

	table := gotabulate.Create(data)
	table.SetHeaders([]string{"FILEPATH", "COMPLEXITY"})
	table.SetAlign("left")

	_, _ = io.WriteString(out, table.Render("grid"))
}
