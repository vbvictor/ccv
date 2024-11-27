package process

import "strings"

type LangName = string
type LangExt = string

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

func GetExtMap(langs []LangName) map[LangExt]struct{} {
	extMap := make(map[LangExt]struct{})
	
	for _, lang := range langs {
		if exts, has := LangMap[lang]; has {
			for _, ext := range strings.Split(exts, ",") {
				extMap[ext] = struct{}{}
			}
		}
	}

	return extMap
}
