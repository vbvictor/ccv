package complexity

type FilesFilterFunc func(files FilesStat) FilesStat

func ApplyFilters(files FilesStat, filters ...FilesFilterFunc) FilesStat {
	result := files

	for _, filter := range filters {
		result = filter(result)
	}

	return result
}

type MinComplexityFilter struct {
	MinComplexity int
}

const (
	MinComplexityDefault = 5
)

func (f MinComplexityFilter) Filter(files FilesStat) FilesStat {
	result := make(FilesStat, 0, len(files))

	for _, file := range files {
		filteredFuncs := make([]FunctionStat, 0)

		for _, fn := range file.Functions {
			if fn.Complexity >= f.MinComplexity {
				filteredFuncs = append(filteredFuncs, fn)
			}
		}

		if len(filteredFuncs) > 0 {
			result = append(result, FileStat{
				Path:      file.Path,
				Functions: filteredFuncs,
			})
		}
	}

	return result
}
