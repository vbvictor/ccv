package complexity

import (
	"fmt"

	"github.com/vbvictor/ccv/pkg/process"
)

type Engine = string

const (
	Gocyclo  Engine = "gocyclo"
	Gocognit Engine = "gocognit"
	Lizard   Engine = "lizard"
)

type Options struct {
	Exts    process.LangExt
	Threads int
	Engine  Engine
	Exclude string
}

var Opts = Options{
	Exts:    "",
	Threads: 1,
	Engine:  Gocyclo,
	Exclude: "",
}

func RunComplexity(repoPath string, opts Options) (FilesStat, error) {
	switch opts.Engine {
	case Gocyclo:
		return RunGocyclo(repoPath, opts)
	case Gocognit:
		return RunGocognit(repoPath, opts)
	case Lizard:
		return RunLizardCmd(repoPath, opts)
	default:
		return nil, fmt.Errorf("unsupported complexity engine: %s", opts.Engine)
	}
}
