package complexity

type FileStat struct {
	Path      string
	Functions FunctionsStat
}

type FilesStat = []FileStat

type FunctionStat struct {
	File       string
	Package    []string
	Name       string
	Line       int
	Length     int
	Complexity int
}

type FunctionsStat = []FunctionStat

type ChurnChunk struct {
	File    string `json:"path"`
	Churn   int    `json:"changes"`
	Added   int    `json:"additions"`
	Removed int    `json:"deletions"`
	Commits int    `json:"commits"`
}

type FileComplexity struct {
	File       string
	Complexity float64
}
