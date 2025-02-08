package complexity

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

type lizardItem struct {
	Name   string `xml:"name,attr"`
	Values []int  `xml:"value"`
}

/*
type lizardFile struct {
	Items []lizardItem `xml:"item"`
}
*/

type lizardMeasure struct {
	Items []lizardItem `xml:"item"`
	Sums  []struct {
		Label string `xml:"label,attr"`
		Value int    `xml:"value,attr"`
	} `xml:"sum"`
	Type string `xml:"type,attr"`
}

type lizardXML struct {
	XMLName  xml.Name        `xml:"cppncss"`
	Measures []lizardMeasure `xml:"measure"`
}

func CheckLizardExecutable() error {
	_, err := exec.LookPath("lizard")
	if err != nil {
		return fmt.Errorf("lizard executable not found, %w", err) // Give pleasant error message
	}

	return nil
}

func RunLizardCmd(repoPath string, opts Options) (FilesStat, error) {
	cmd := []string{"lizard", "-s", "cyclomatic_complexity", "-m", "-X", "-t", strconv.Itoa(opts.Threads)}

	if opts.Exts != "" {
		allowedExts := strings.Split(opts.Exts, ",")
		for _, ext := range allowedExts {
			cmd = append(cmd, "-l", ext)
		}
	}

	if opts.Exclude != "" {
		cmd = append(cmd, "--exclude", opts.Exclude)
	}

	cmd = append(cmd, repoPath)

	lizardCmd := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec // lizard is a trusted binary

	output, err := lizardCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run lizard: %w", err)
	}

	lizard, err := readLizardXML(bytes.NewReader(output))
	if err != nil {
		return nil, fmt.Errorf("failed to read lizard output: %w", err)
	}

	return parseLizard(lizard)
}

var (
	lizardFileRe       = regexp.MustCompile(`(.*?)\s+at\s+(.*?):(\d+)`)
	errLizard          = errors.New("lizard error")
	expectedMatchCount = 4
)

func parseItem(item lizardItem) (FunctionStat, error) {
	matches := lizardFileRe.FindStringSubmatch(item.Name)

	if len(matches) != expectedMatchCount {
		return FunctionStat{}, fmt.Errorf("invalid function format %s, %w", item.Name, errLizard)
	}

	funcParts := strings.Split(matches[1], "::")
	name := funcParts[len(funcParts)-1]

	if idx := strings.Index(name, "("); idx != -1 {
		name = name[:idx]
	}

	lineNum, _ := strconv.ParseInt(matches[3], 10, 32)

	return FunctionStat{
		File:       matches[2],
		Name:       name,
		Package:    funcParts[:len(funcParts)-1],
		Line:       int(lineNum),
		Length:     item.Values[1],
		Complexity: item.Values[2],
	}, nil
}

func readLizardXML(r io.Reader) (*lizardXML, error) {
	var lizard lizardXML
	if err := xml.NewDecoder(r).Decode(&lizard); err != nil {
		return nil, fmt.Errorf("failed to decode XML, %w", err)
	}

	return &lizard, nil
}

func parseLizard(lizard *lizardXML) (FilesStat, error) {
	type Filename = string

	fileStat := make(map[Filename]FileStat)

	filesIdx := slices.IndexFunc(lizard.Measures, func(l lizardMeasure) bool {
		return l.Type == "File"
	})

	for _, file := range lizard.Measures[filesIdx].Items {
		fileStat[file.Name] = FileStat{Path: file.Name}
	}

	funcIdx := slices.IndexFunc(lizard.Measures, func(l lizardMeasure) bool {
		return l.Type == "Function"
	})

	for _, function := range lizard.Measures[funcIdx].Items {
		funcStat, err := parseItem(function)
		if err != nil {
			return nil, err
		}

		file := fileStat[funcStat.File]
		file.Functions = append(fileStat[funcStat.File].Functions, funcStat)
		fileStat[funcStat.File] = file
	}

	return maps.Values(fileStat), nil
}
