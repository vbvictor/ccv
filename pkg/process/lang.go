package process

type (
	LangName = string
	LangExt  = string
)

var LangMap = map[LangName]LangExt{
	"cpp":        "h,hpp,c,cpp",
	"csharp":     "cs",
	"go":         "go",
	"java":       "java",
	"javascript": "js",
	"typescript": "ts",
	"python":     "py",
	"ruby":       "rb",
	"rust":       "rs",
	"php":        "php",
}

func GetExtMap(exts []LangExt) map[LangExt]struct{} {
	if len(exts) == 0 {
		return nil
	}

	extMap := make(map[LangExt]struct{})

	for _, ext := range exts {
		extMap[ext] = struct{}{}
	}

	return extMap
}
