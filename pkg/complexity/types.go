package complexity

type FileStat struct {
	Path      string
	Functions FunctionsStat
}

type FilesStat = []*FileStat

type FunctionStat struct {
	File      string
	Package   []string
	Name      string
	Line      uint
	Length    uint
	Compexity uint
}

type FunctionsStat = []FunctionStat

type ChurnChunk struct {
	File    string `json:"path"`
	Churn   uint   `json:"changes"`
	Added   uint   `json:"additions"`
	Removed uint   `json:"deletions"`
	Commits uint   `json:"commits"`
}
