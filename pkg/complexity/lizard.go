package complexity

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Engine = string

var (
	Lizard Engine = "lizard"
)

type ComplexityOptions struct {
	Extensions string
	Threads    int
}

var ComplexityOpts = ComplexityOptions{
	Extensions: "",
	Threads: 1,
}

type lizardItem struct {
	Name   string `xml:"name,attr"`
	Values []int  `xml:"value"`
}

type lizardFile struct {
	Items []lizardItem `xml:"item"`
}

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

var lizardFileRe = regexp.MustCompile(`(.*?)\s+at\s+(.*?):(\d+)`)

func RunLizardCmd(repoPath string, opts ComplexityOptions) (FilesStat, error) {
	_, err := exec.LookPath("lizard")
	if err != nil {
		return nil, err
	}

	cmd := []string{"lizard", "-s", "cyclomatic_complexity", "-m", "-X", "-t", strconv.Itoa(opts.Threads)}

	if opts.Extensions != "" {
		// Walk filepath to be sure such files exist
		// If no files exist print that no files and return empty FilesStat
		allowedExts := strings.Split(opts.Extensions, ",")
		for _, ext := range allowedExts {
			cmd = append(cmd, "-l", ext)
		}
	}

	cmd = append(cmd, repoPath)
	lizardCmd := exec.Command(cmd[0], cmd[1:]...)
	output, err := lizardCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run lizard: %w", err)
	}

	lizard, err := ReadLizardXML(bytes.NewReader(output))
	if err != nil {
		return nil, fmt.Errorf("failed to parse lizard output: %w", err)
	}

	return ParseLizard(lizard)
}

func parseItem(item lizardItem) (FunctionStat, error) {
	re := regexp.MustCompile(`(.*?)\s+at\s+(.*?):(\d+)`)
	matches := re.FindStringSubmatch(item.Name)

	if len(matches) != 4 {
		return FunctionStat{}, fmt.Errorf("invalid function format: %s", item.Name)
	}

	funcParts := strings.Split(matches[1], "::")
	name := funcParts[len(funcParts)-1]

	if idx := strings.Index(name, "("); idx != -1 {
		name = name[:idx]
	}

	lineNum, _ := strconv.ParseUint(matches[3], 10, 32)

	return FunctionStat{
		File:      matches[2],
		Name:      name,
		Package:   funcParts[:len(funcParts)-1],
		Line:      uint(lineNum),
		Length:    uint(item.Values[1]),
		Compexity: uint(item.Values[2]),
	}, nil
}

func ReadLizardXML(r io.Reader) (*lizardXML, error) {
	var lizard lizardXML
	if err := xml.NewDecoder(r).Decode(&lizard); err != nil {
		return nil, err
	}

	return &lizard, nil
}

func ParseLizard(lizard *lizardXML) (FilesStat, error) {
	type Filename = string
	fileMap := make(map[Filename]*FileStat)

	filesIdx := slices.IndexFunc(lizard.Measures, func(l lizardMeasure) bool {
		return l.Type == "File"
	})

	for _, file := range lizard.Measures[filesIdx].Items {
		fileMap[file.Name] = new(FileStat)
		fileMap[file.Name].Path = file.Name
	}

	funcIdx := slices.IndexFunc(lizard.Measures, func(l lizardMeasure) bool {
		return l.Type == "Function"
	})

	for _, function := range lizard.Measures[funcIdx].Items {
		stat, err := parseItem(function)
		if err != nil {
			return nil, err
		}

		fileMap[stat.File].Functions = append(fileMap[stat.File].Functions, stat)
	}

	// Convert map to slice
	result := make([]*FileStat, 0, len(fileMap))
	for _, fileStat := range fileMap {
		result = append(result, fileStat)
	}

	return result, nil
}
