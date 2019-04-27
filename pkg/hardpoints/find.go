package hardpoints

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var hardpointRegex = regexp.MustCompile(`HARDPOINT\:(?P<tags>(\s*\w+\s*,?)+)$`)

var ignoreList = map[string]bool{
	".git":         true,
	".svn":         true,
	"vendor":       true,
	"node_modules": true,
}

// TODO: Make parallel over each file
func FindHardpoints(root string) (hps []*Point) {
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && ignoreList[info.Name()] {
			return filepath.SkipDir
		}
		if nhps := extractHardpoints(path); nhps != nil {
			hps = append(hps, nhps...)
		}
		return err
	})
	return hps
}

// TODO: Make parallel over each line in the text
func extractHardpoints(path string) (out []*Point) {
	f, err := os.Open(path)
	if err != nil {
		return out
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var line int
	for scanner.Scan() {
		line++

		hp := maybeExtractHardpoint(scanner.Text())
		if hp == nil {
			continue
		}
		hp.Line = line
		hp.Path = path
		out = append(out, hp)
	}

	return out
}

func maybeExtractHardpoint(s string) *Point {
	parts := make(map[string]string)

	found := hardpointRegex.FindStringSubmatch(s)
	if found == nil {
		return nil
	}

	for i, name := range hardpointRegex.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		parts[name] = found[i]
	}

	return &Point{
		Tags: extractTags(parts["tags"]),
	}
}

func extractTags(s string) (out []string) {
	for _, v := range strings.Split(s, ",") {
		out = append(out, strings.TrimSpace(v))
	}
	return out
}
