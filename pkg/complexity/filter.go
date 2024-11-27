package complexity

type FilesFilterFunc func(files FilesStat) FilesStat

func ApplyFilters(files FilesStat, filters ...FilesFilterFunc) FilesStat {
	result := files

	for _, filter := range filters {
		result = filter(result)
	}

	return result
}

type ComplexityFilter struct {
	MinComplexity uint
}

func (f ComplexityFilter) Filter(files FilesStat) FilesStat {
	result := make(FilesStat, 0, len(files))

	for _, file := range files {
		filteredFuncs := make([]FunctionStat, 0)
		for _, fn := range file.Functions {
			if fn.Compexity >= f.MinComplexity {
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
